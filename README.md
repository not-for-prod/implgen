# Implgen

Install:

```bash
go install github.com/not-for-prod/implgen@latest
```

## Basic

> i wanted impementation generator that is working like generator inside IDE but cli

Flags:

- `src` - source file filepath
- `dst` - destination path
- `interface-name` - specify which interface u need to implement
- `with-otel` - if true and first method argument is from `context` package will generate

    ```go
    ctx, span := otel.Tracer("").Start(ctx, "AbobaImplementation.Create")
	defer span.End()
    ```

HOW I USE IT

```bash
implgen --src ./example/in/aboba.go --dst ./example/out/basic --interface-name Aboba --with-otel
```

## Repo

> i wanted impementation generator for repo layer compatable with `[github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx)` and 
> `[github.com/avito-tech/go-transaction-manager/sqlx](https://github.com/avito-tech/go-transaction-manager/sqlx)`
> that will generate basic stuff aka `.sql` files method files with basic stuff i m tired to fill

Flags:

- `src` - source file filepath
- `dst` - destination path
- `interface-name` - specify which interface u need to implement

HOW I USE IT

```bash
implgen repo --src ./example/in/aboba.go --dst ./example/out/repo --interface-name AbobaRepository
```

## How it was made

As base i used [github.com/golang/mock/blob/main/mockgen](https://github.com/golang/mock/blob/main/mockgen) but as it was private i copied it and got it's guts out.
Than did not enjoyed how generation was made so i took [google.golang.org/protobuf/compiler/protogen](https://google.golang.org/protobuf/compiler/protogen) generator.
And after that had problems with imports so i used [golang.org/x/tools/imports](https://golang.org/x/tools/imports) and `go/format` on the top.
More features like span generation and that's it.