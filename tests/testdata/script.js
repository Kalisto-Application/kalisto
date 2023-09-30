firstRes = SequenceService.First({
  value: 0,
  rpc: '',
});

secondRes = SequenceService.Second({
  value: firstRes.value,
  rpc: firstRes.rpc,
});

thirdRes = SequenceService.Third({
  value: secondRes.value,
  rpc: secondRes.rpc,
});
