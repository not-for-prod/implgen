# Implgen

Install:

```bash
go install github.com/not-for-prod/implgen@latest
```

> i wanted impementation generator that is working like generator inside IDE but cli

As base i used [github.com/golang/mock/blob/main/mockgen](https://github.com/golang/mock/blob/main/mockgen) but as it was private i copied it and got it's guts out. 
Than did not enjoyed how generation was made so i took [google.golang.org/protobuf/compiler/protogen](https://google.golang.org/protobuf/compiler/protogen) generator.
And after that had problems with imports so i used [golang.org/x/tools/imports](https://golang.org/x/tools/imports) and `go/format` on the top.
More features like span generation and that's it. 

Few moments later added [golden linter](https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322) and modified it with:

```yaml
  exclude-files:
    - internal/mockgen/
    - internal/implgen/mockgen.go
```

Flags:

- `src` - source file filepath
- `dst` - destination path
- `with-otel` - if true and first method argument is from `context` package will generate 

    ```go
    ctx, span := otel.Tracer("").Start(ctx, "AbobaImplementation.Create")
	defer span.End()
    ```

HOW I USE IT

```bash
implgen --src ./example/in/aboba.go --dst ./example/out --with-otel
```