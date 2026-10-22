package assert

var (
    EnablePanic = true   //panic开关
    //WriteFail 写日志
    WriteFail = func(template string, args ...interface{}) {}
)

//doFail 失败处理
func doFail(msg string) {
    if EnablePanic {
        panic(msg)
    }
    WriteFail("[ASSERT FAIL] %s", msg)
}
