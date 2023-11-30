// go:build ignore

#include "xdp_bandwidth.c"
#include "xdp_packetloss.c"

char _license[] SEC("license") = "GPL";

SEC("xdp")
int xdp_main(struct xdp_md *ctx)
{
  int action = xdp_packetloss(ctx);
  if (action != XDP_PASS)
  {
    return action;
  }
  return xdp_bandwidth_limit(ctx);
}