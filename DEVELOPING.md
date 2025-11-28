# Developing

## TFLint SDK Parsing Methods

The TFLint SDK provides two primary methods for inspecting Terraform configuration: Module Content Mode and AST Mode. We strictly utilize AST Mode in this ruleset.

### Module Content Mode (`runner.GetModuleContent`)

This is the higher-level abstraction provided by the SDK. It evaluates the Terraform configuration against a schema (often `hclext`).

*   **Pros:**
    *   Simplifies access to attribute values (e.g., getting a string value directly).
    *   Handles variable interpolation and locals evaluation to some extent.
*   **Cons:**
    *   **Pre-evaluation:** The SDK pre-parses the source and evaluates logic like `count` and `for_each`. If a resource has `count = 0`, it is effectively removed from the content returned to the rule. This makes it impossible to lint the *configuration* of that resource if it happens to be disabled in the current context.
    *   Hides raw syntax details (comments, specific token locations).

### AST Mode (`runner.GetFiles`)

This method provides access to the raw Abstract Syntax Tree (AST) of the HCL files.

*   **Pros:**
    *   **Complete Visibility:** You see the code exactly as written. Resources with `count = 0` are present.
    *   Access to all tokens, comments, and exact source ranges.
*   **Cons:**
    *   **Complexity:** Requires manual traversal of the AST (`hclsyntax` package).
    *   No automatic evaluation of expressions or variables. You are working with raw expressions.

### Current Usage

We strictly use **AST Mode** for all rules in this ruleset.

Our primary concern is the *style* of the code as written, not the evaluation of the graph at plan/apply time. Using AST mode ensures that we lint all code, including resources that might be planned away explicitly (e.g., `count = 0`) or conditionally by Terraform's evaluation logic. This provides a consistent linting experience regardless of the current variable inputs.

Any validation of attribute *values* (which AST mode makes difficult) should be delegated to Terraform's native `validation`, `precondition`, or `check` blocks.
