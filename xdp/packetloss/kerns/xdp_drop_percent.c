// go:build ignore
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <string.h>

#define MAX_MAP_ENTRIES 1
struct
{
  __uint(type, BPF_MAP_TYPE_LRU_HASH);
  __type(key, __u32);
  __type(value, __s32);
  __uint(max_entries, MAX_MAP_ENTRIES);
} drop_rate_map SEC(".maps");

SEC("xdp")
int xdp_drop_percent(struct xdp_md *ctx)
{
  // the map is looked up only once for performance purposes
  static __s32 drop_rate = -1;
  if (drop_rate == -1)
  {
    __u32 key = 0;
    __s32 *drop_rate_ptr = bpf_map_lookup_elem(&drop_rate_map, &key);
    if (!drop_rate_ptr)
      return XDP_ABORTED;

    drop_rate = *drop_rate_ptr;
  }

  if (drop_rate > 0 && bpf_get_prandom_u32() % 100 < drop_rate)
  {
    return XDP_DROP;
  }

  return XDP_PASS;
}

char _license[] SEC("license") = "GPL";