res = BookStore.GetBook({
  id: 'book-id',
  double: 3.14,
  float: 4.14,
  int32: 1,
  int64: 2,
  uint32: 3,
  uint64: 4,
  sint32: 5,
  sint64: 6,
  fixed32: 7,
  fixed64: 8,
  sfixed32: 9,
  sfixed64: 10,
  bool: true,
  bytes: '{json: true}',
  enum: 3,
  Book: {
    id: 'another-id',
  },
  strToStr: { str: 'to-str' },
  intToBooks: {
    14: {
      id: 'book-id',
    },
  },
  btoe: { true: 3 },
  strings: ['str1', 'str2'],
  enums: [4],
  uints: [111],
  books: [
    {
      id: 'id-in-slice',
    },
  ],
  deepNestedBook: {
    hasNested: {
      deepNested: {
        hasNested: {},
      },
    },
  },
  repeatedNestedBook: [
    {
      hasNested: {},
    },
  ],
  someBook: { oneId: 'toAnother' },
  // someBook: {"oneEnum": 0},
  anotherBook: {
    anotherBookObject: {
      id: 'anotherObjectId',
    },
  },
  // anotherBook: {"anotherNestedBook": {
  //     hasNested: {
  //     },
  //   }},
  // 8760h0m0s or 1 year
  dur: 31536000000000000,
  // 2022-04-04T00:00:00.000Z or 1649030400000 in ms
  time: new Date(Date.UTC(2022, 3, 4, 0, 0, 0, 0)),
  // unused field
  time2: Date.now(),
  l1: {
    f1: {
      f1: {
        f1: '1',
        f2: '2',
      },
      f2: {
        f1: '3',
        f2: '4',
      },
    },
    f2: {
      f1: {
        f1: '5',
        f2: '6',
      },
      f2: {
        f1: '7',
        f2: '8',
      },
    },
  },
  l2: {
    f1: {
      f1: {
        f1: '11',
        f2: '22',
      },
      f2: {
        f1: '33',
        f2: '44',
      },
    },
    f2: {
      f1: {
        f1: '55',
        f2: '66',
      },
      f2: {
        f1: '77',
        f2: '88',
      },
    },
  },
});

res = BookStore.GetBook({
  id: res.id,
  double: res.double,
  float: res.float,
  int32: res.int32,
  int64: res.int64,
  uint32: res.uint32,
  uint64: res.uint64,
  sint32: res.sint32,
  sint64: res.sint64,
  fixed32: res.fixed32,
  fixed64: res.fixed64,
  sfixed32: res.sfixed32,
  sfixed64: res.sfixed64,
  bool: res.bool,
  bytes: res.bytes,
  enum: res.enum,
  Book: res.Book,
  strToStr: res.strToStr,
  intToBooks: res.intToBooks,
  btoe: res.btoe,
  strings: res.strings,
  enums: res.enums,
  uints: res.uints,
  books: res.books,
  deepNestedBook: res.deepNestedBook,
  repeatedNestedBook: res.repeatedNestedBook,
  someBook: res.someBook,
  anotherBook: res.anotherBook,
  // 8760h0m0s or 1 year
  dur: res.dur,
  // 2022-04-04T00:00:00.000Z or 1649030400000 in ms
  time: res.time,
  // unused field
  l1: res.l1,
  l2: res.l2,
});
