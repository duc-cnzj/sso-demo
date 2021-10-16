const FixiOSBug = require("./FixiOSBug")
// @ponicode
describe("FixiOSBug.default.mounted", () => {
    test("0", () => {
        let callFunction = () => {
            FixiOSBug.default.mounted()
        }
    
        expect(callFunction).not.toThrow()
    })
})
