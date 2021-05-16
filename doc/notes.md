# Notes

## Spread restrictions

Consider using a filter on ... lifts to allow lifting only part of the scope.

Example:
```dg
  ...Wasm[inst]
```

which could be sugared as,

```dg
  from Wasm import [inst]
```

splat or wildcard imports leads to magic identiifers (e.g. where is this
identifer defined) in scope and can lead to change at a distance problems as
the scope iported evolves.

Consider not allowing wildcard imports at all (or discourage them as Go does or
Google Java lint rules do).

Consider making this general for all spread uses.

Downside:
  - spread in a function parameter list allows higher order function
    declarations as the function evolves

    ```dg
      fun Button(text: string, click: { -> Unit }) { 
          ...
       }

      fun BancyButton(fancy: Fancy, ...Button.parameters) {
          ...
          Button(...button.parameters)
      }
    ```

    which could be desugared into,

    ```
    let Button = { text: String, click { -> Unit } -> 
        ...
    }
    let FancyButton = { fancy: Fancy, text: string, click: { -> Unit } ->
        ... use fancy ...
        Button(texxt: text, click: click)
    }
    ```

    Requring explicit qualification here would reduce the utility of this 
    technique.

    However, the qualification could be used to identify important lifted
    parameters. In the above the parameters are passed anonymously, lifting
    doesn't lift the symbols, just the the parameters. Any lifted conflicting
    parameters are alpha renamed. For example,

    ```
    fun FancyButton(text: String, fancy: Fancy, ...Button.parameters) {
        ... use fancy ...
        Button(text, ...Button.parameter)
    }
    ```

    could desugar into,

    ```
    let Button = { text: String, fancy: Fancy, `Button.text`: string, click: { -> Unit} -> 
        ... use fancy ...
        Button(text: text, click: click)
    ```

    The lifted alpha renamed `Button.text` is unfortunate, however. Consider
    eliding `Button.text` or allowing it to be filtered as such as,

    ```
    fun FancyButton(text: String, fancy: Fancy, ...Button.parameters[-text]) {
        ...
    }
    ```

  - spread at a record level is similar to a function parameters (due to
    correpondence of record to parameters and record to result) but the raw
    layout doesn't have a type impact as it does in parameters as the
    existential (private) fields don't need to be lifted into the record type.

## Type levels

Dyego0 is an nominal typed system without polymorphic references but can have witnessed references (e.g. interfaces or protocols) and record embedding an lifting (similar to Go but wihtout polymorphic pointers). The intent is that DyegoN (or just Dyego) will introduce desugaring layers that allow more powerful type systems in the layer such as System Fω, or a matching based variant. The intent is that type system would be added as a library layer on a more primitive language. This hopefully allow for experimentation and as well as deprecation in the type system itself. That is a module can import a type system and the type systems do not need to agree across module boundaries as long as they can be desugared into a lower level type system they do agree on. The lowest level being Dyego0.

