# Implgen

> i wanted impementation generator that is working like generator inside IDE but cli

As base i used [github.com/golang/mock/blob/main/mockgen](https://github.com/golang/mock/blob/main/mockgen) but as it was private i copied it and got it's guts out. 
Than did not enjoyed how generation was made so i took [google.golang.org/protobuf/compiler/protogen](https://google.golang.org/protobuf/compiler/protogen) generator.
And after that had problems with imports so i used [golang.org/x/tools/imports](https://golang.org/x/tools/imports) and `go/format` on the top.
More features like span generation and that's it.

Flags:

- `src` - source file filepath
- `dst` - destination path
- `with-otel` - if true and method has ctx will generate 

    ```go
    ctx, span := otel.Tracer("").Start(ctx, "AbobaImplementation.Create")
	defer span.End()
    ```

HOW I USE IT

```bash
implgen --src ./example/in/aboba.go --dst ./example/out --with-otel
```

Also i left 