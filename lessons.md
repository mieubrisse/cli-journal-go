Lots of tricky Go stuff going on here to know about:
- The "Update" command is intended to return a **sub-instance of the tea.Model**, not the actual tea.Model (i.e., it doesn't actually conform to the `Model` interface)
- Buuble Tea wants stuff passed by-value (especially with the `Update` command) so the idea is that you have to re-assign after you call `Update` on a sub-component
    - This is REALLY confusing (how do you have a single model element that's used by multiple componenets?)
- It's weird that you can only update a model in a single location
