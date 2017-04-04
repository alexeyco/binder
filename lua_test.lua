-- see lua_test.go

local r = require('reverse')

assert(r.string('ABCDE') == 'EDCBA', 'ABCDE != EDCBA')
assert(r.string('01234') == '43210', '01234 != 43210')

assert(r._STRING == 'foobar', 'Wrong string field value')
assert(r._NUMBER == 123, 'Wrong number field value')
assert(r._BOOL, 'Wrong bool field value')
