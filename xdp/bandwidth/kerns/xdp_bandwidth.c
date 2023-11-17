// go:build ignore
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <string.h>

#define MAX_MAP_ENTRIES 1
#define NANOS_PER_SEC 1000000000UL
// A time window is used to measure the average bandwidth.
#define TIME_WINDOW_SEC 5

struct
{
  __uint(type, BPF_MAP_TYPE_HASH);
  __type(key, __u32);
  __type(value, __u64);
  __uint(max_entries, MAX_MAP_ENTRIES);
} last_packet_timestamp SEC(".maps");

struct
{
  __uint(type, BPF_MAP_TYPE_LRU_HASH);
  __type(key, __u32);
  __type(value, __u64);
  __uint(max_entries, MAX_MAP_ENTRIES);
} byte_counter SEC(".maps");

struct
{
  __uint(type, BPF_MAP_TYPE_LRU_HASH);
  __type(key, __u32);
  __type(value, __u64); // Bytes per second
  __uint(max_entries, MAX_MAP_ENTRIES);
} bandwidth_limit_map SEC(".maps");

SEC("xdp")
int xdp_bandwidth_limit(struct xdp_md *ctx)
{
  __u64 current_timestamp = bpf_ktime_get_ns();
  __u32 key = 0;
  __u64 *last_time_window_start = bpf_map_lookup_elem(&last_packet_timestamp, &key);
  if (!last_time_window_start)
  {
    bpf_map_update_elem(&last_packet_timestamp, &key, &current_timestamp, BPF_ANY);
    return XDP_PASS;
  }

  __u64 packet_size = (__u64)(ctx->data_end - ctx->data);

  __u64 *byte_count_ptr = bpf_map_lookup_elem(&byte_counter, &key);
  if (!byte_count_ptr)
  {
    bpf_map_update_elem(&byte_counter, &key, &packet_size, BPF_ANY);
  }
  else
  {
    *byte_count_ptr += packet_size;
  }

  static __u64 time_window_ns = TIME_WINDOW_SEC * NANOS_PER_SEC;
  if (current_timestamp - *last_time_window_start >= time_window_ns)
  {
    // Reset byte counter and start of the next time window
    __u64 reset_value = 0;
    bpf_map_update_elem(&byte_counter, &key, &reset_value, BPF_ANY);
    bpf_map_update_elem(&last_packet_timestamp, &key, &current_timestamp, BPF_ANY);
  }

  // Look it up only once for performance purposes
  static __u64 allowed_bytes = 0; // number of bytes per window
  if (allowed_bytes == 0)
  {
    __u64 *bandwidth_limit_ptr = bpf_map_lookup_elem(&bandwidth_limit_map, &key);
    if (!bandwidth_limit_ptr)
      return XDP_ABORTED;

    allowed_bytes = (*bandwidth_limit_ptr / 8 * TIME_WINDOW_SEC); // divide by 8 to convert from bits to bytes
  }

  __u64 *accumulated_bytes = bpf_map_lookup_elem(&byte_counter, &key);
  if (!accumulated_bytes)
    return XDP_ABORTED;

  if (*accumulated_bytes > allowed_bytes)
    return XDP_DROP;

  return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
