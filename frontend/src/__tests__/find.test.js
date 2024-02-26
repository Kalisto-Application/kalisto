/* eslint-disable */
function find(arr, callback) {

    for (let index = 0; index < arr.length; index++) {
        if (callback(arr[index])) {
            return arr[index]
        }
    }


}


test('find when 2', () => {
    const given = [1, 2, 3]
    const callback = (num) => {
        return num === 2;
    }
    const result = find(given, callback)

    expect(result).toBe(2)
})


test('find when even', () => {
    const given = [1, 2, 3]
    const callback = (num) => {
        return num % 2 === 0;
    }
    const result = find(given, callback)

    expect(result).toBe(2)
})


test('find first', () => {
    const given = [1, 2, 3]
    const callback = (num) => {
        return true;
    }
    const result = find(given, callback)

    expect(result).toBe(1)
})

test('find doubledigit', () => {
    const given = [1, 11, 999]
    const callback = (num) => {
        return num >= 10;
    }
    const result = find(given, callback)

    expect(result).toBe(11)
})

// testCases = [
//     {
//         given: [1, 2, 3],
//         callback: (num) => {
//             return num === 2;
//         },
//         expected: 2,
//         name: 'when 2'
//     }
// ]

// testCases.forEach(case => {
//     test(case.name, () => {
//     const result = find(case.given, case.callback)

// expect(result).toBe(case.expected)
//     })
// })