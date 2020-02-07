# syncs

### RefCounter

RefCounter manually manages reference count of a object which is created with RefCounter instance, then calls release function only if the reference count reaches zero.

Recommends call Clone() method of RefCounter whenever RefCounter instance is handed over by function argument, then defer Release() method at first in that function.

The usual use pattern is like below

```go
	rc := syncs.NewRefCounter(&SomeObject{}, func(obj interface{}) {
		obj.(io.Closer).Close()
    })
    defer rc.Release()

    // in some other block
    func(rc RefCounter){
        defer rc.Release()

        object := rc.Interface().(*SomeObject)

        // do some work with object
        object.DoSomeWork()

    }(rc.Clone())

```

```plantuml
Client -> RefCounter **: rc = New(object, release_func)
group in another block of code
Client -> RefCounter: clonedRc = rc.Clone()
Client -> SomeElse ++: do_something(clonedRc)
SomeElse -> RefCounter: object = Interface()
ref over SomeElse
do some work with object
end
SomeElse -> RefCounter: Release()
return
Client -> RefCounter: Release()
end
```
