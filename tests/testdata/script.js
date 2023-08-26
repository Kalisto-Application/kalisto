firstRes = SequenceService.First({
    value: 0,
    rpc: '',
})

secondRes = SequenceService.Second(firstRes)

thirdRes = SequenceService.Third(secondRes)
