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
} packetloss_rate_map SEC(".maps");

int xdp_packetloss(struct xdp_md *ctx)
{
  __u32 key = 0;
  __s32 *drop_rate_ptr = bpf_map_lookup_elem(&packetloss_rate_map, &key);
  if (!drop_rate_ptr)
  {
    // if it has not set by the user space program,
    // or the service is not started yet
    return XDP_PASS;
  }

  if (*drop_rate_ptr == 0)
  {
    // If the service is stopped
    return XDP_PASS;
  }

  if (bpf_get_prandom_u32() % 100 < *drop_rate_ptr)
  {
    return XDP_DROP;
  }

  return XDP_PASS;
}