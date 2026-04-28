// ---------------- Pre-made arrays ----------------

tabledata = [
    ["h=t&sid=", "6="],
    ["objects=", "1="],
    ["sid=", "2="],
    ["&ox=0&types=0&fe=0", "3=3"],
    ["&types=0&fe=0", "5=5"],
    ["&h=t&p=", "4="]
];

tabledataspecial = [
    ["=", "ZZZX1"],
    ["&", "ZZZX2"],
    [",", "ZZZX3"],
    [".", "ZZZX4"],
    [" ", "ZZZX5"],
    ["-", "ZZZX6"],
    ["/", "ZZZX7"],
    ["%", "ZZZX8"]
];

pairs = [
    ["=", "Q"],
    ["&", "Z"],
    [",", "X"],
    [".", "Y"],
    [" ", "V"],
    ["-", "W"]
]

pattern = [4, 22, 5, 37, 26, 17, 33, 15, 39, 11, 45, 20, 2, 40, 19, 36, 28, 38, 30, 41, 44, 42, 7, 24, 14, 27, 35, 25, 12, 1, 43, 23, 6, 16, 3, 9, 47, 46, 48, 50, 21, 10, 49, 32, 18, 31, 29, 34, 13, 8];

// ---------------- Helper functions ----------------

function isEmpty(str) {
    return (!str || 0 === str.length);
}

function toString(str) {
    if (isEmpty(str)) return "";
    return "" + str;
}

// ---------------- Scramble Functinons ----------------

function tablespecial(result, reverse=false) {
    for (let index = 0; index < tabledataspecial.length; index++) {
        let key = tabledataspecial[index];
        if (reverse) result = result.replaceAll(key[1], key[0]);
        else result = result.replaceAll(key[0], key[1]);
    }
    return result;
}

function tableshort(result, reverse=false) {
    for (let index = 0; index < tabledata.length; index++) {
        let key = tabledata[index];
        if (reverse) result = result.replaceAll(key[1], key[0]);
        else result = result.replaceAll(key[0], key[1]);
    }
    return result;
}

function modKey(ch, reverse=false) {
    if (ch >= 97 && ch <= 122) {
        if (reverse) return (97 + (ch - 106 + 26) % 26);
        else return (97 + (ch - 88) % 26);
    }
    if (ch >= 49 && ch <= 57) {
        if (reverse) return (49 + (ch - 53 + 9) % 9);
        else return (49 + (ch - 45) % 9);
    }
    return ch;
}

function scrambleChar(ch, reverse=false) {
    for (let index = 0; index < pairs.length; index++) {
        let pair = pairs[index];
        if (ch == pair[0]) return pair[1];
        if (ch == pair[1]) return pair[0];
    }return String.fromCharCode(modKey(ch.charCodeAt(0), reverse));
}

function swap(result, from, to) {
    if ((from < 0) || (from >= result.length)) return;
    if ((to < 0) || (to >= result.length)) return;
    let fromChar = result[from];
    result[from] = result[to];
    result[to] = fromChar;
}

function swapPattern(result, reverse=false) {
    let steps = Math.ceil(result.length / pattern.length);
    for (let step = 0; step < steps; step++) {
        for (let index = 1; index < pattern.length; index += 2) {
            if (reverse) {
                swap(result, pattern[index - 1] + step * pattern.length, pattern[index] + step * pattern.length);
            } else {
                swap(result, pattern[index] + step * pattern.length, pattern[index - 1] + step * pattern.length);
            }
        }
    }
}

function swapChar(result, reverse=false) {
    let split = result.split("");
    for (let index = 0; index < split.length; index++) {
        split[index] = scrambleChar(split[index], reverse);
    }
    swapPattern(split, reverse);
    return split.join("");
}

function scramble(query, reverse=false) {
    if (isEmpty(query)) return query;
    if (query.length < 2) return query;
    if (query.substring(0, 2) === "i=") return query;

    let result = decodeURIComponent(query);

    let steps = [
        tableshort,
        swapChar,
        tablespecial
    ];

    let orderedSteps = steps.slice();
    if (reverse) orderedSteps.reverse();

    for (let f of orderedSteps) {
        result = f(result, reverse);
    }

    return encodeURIComponent(result);
}

let keyValues = "sid=3&objects=957.8&p=0.m,20260607.x";

console.log(scramble(keyValues))

// TKDAT-1 + ZBASS-1
//console.log(decodeURIComponent(scramble("14QYZY3Q7Z09Q2v04Y62Q35y90Z6g19x1Y60gQXX70", true)));
// h=t&sid=3&objects=957.8,932.8&ox=0&p=0.m,20260607.x&types=0&fe=0

// TKDAT-1
//console.log(decodeURIComponent(scramble("6y7YQQu1QZnZQ28Z48beZ35690Q", true)));
// sid=3&objects=957.8&p=4&l=sv&e=2604

/*

Params : 
- sid (school id?)
- objects (classes/courses)
- p (period)
- l (language)
- e (?)

- h (?)
- ox (?)
- types (?)
- fe (?)

*/

// ---------------- Query builder ----------------

