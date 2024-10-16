package testdata

import (
    orm "geektime-go/day5_orm/my_orm_mysql"
    
    "database/sql"
    
    sqlx "database/sql"
    
)




const (
    
        UserName = "User"
    
        UserAge = "User"
    
        UserNickName = "User"
    
        UserPhone = "User"
    
        UserPicture = "User"
    

)

    
        func UserNameLt(val string) orm.Predicate {
        return orm.C("Name").Lt(val)
        }
    
        func UserNameEq(val string) orm.Predicate {
        return orm.C("Name").Eq(val)
        }
    

    
        func UserAgeLt(val int) orm.Predicate {
        return orm.C("Age").Lt(val)
        }
    
        func UserAgeEq(val int) orm.Predicate {
        return orm.C("Age").Eq(val)
        }
    

    
        func UserNickNameLt(val *sql.NullString) orm.Predicate {
        return orm.C("NickName").Lt(val)
        }
    
        func UserNickNameEq(val *sql.NullString) orm.Predicate {
        return orm.C("NickName").Eq(val)
        }
    

    
        func UserPhoneLt(val *sqlx.NullString) orm.Predicate {
        return orm.C("Phone").Lt(val)
        }
    
        func UserPhoneEq(val *sqlx.NullString) orm.Predicate {
        return orm.C("Phone").Eq(val)
        }
    

    
        func UserPictureLt(val []byte) orm.Predicate {
        return orm.C("Picture").Lt(val)
        }
    
        func UserPictureEq(val []byte) orm.Predicate {
        return orm.C("Picture").Eq(val)
        }
    


const (
    
        UserDetailAddress = "UserDetail"
    

)

    
        func UserDetailAddressLt(val string) orm.Predicate {
        return orm.C("Address").Lt(val)
        }
    
        func UserDetailAddressEq(val string) orm.Predicate {
        return orm.C("Address").Eq(val)
        }
    

