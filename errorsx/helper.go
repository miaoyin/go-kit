package errorsx

import "errors"

//JoinErrors 合并error
func JoinErrors(errs []error) error {
    var list []error
    for _, e := range errs {
        if e != nil {
            list = append(list, e)
        }
    }
    if len(list) == 0 {
        return nil
    }
    return errors.Join(list...)
}
