firstRes = SequenceService.First({
  value: 0,
  rpc: '',
});

secondRes = SequenceService.Second({
  value: firstRes.body.value,
  rpc: firstRes.body.rpc,
});

thirdRes = SequenceService.Third({
  value: secondRes.body.value,
  rpc: secondRes.body.rpc,
});
