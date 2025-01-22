# vim / run sql

When I run `<Leader>r` from a `.sql` file in Vim,
the file's contents are run in my project's Postgres database through `psql`
and the output is sent to a Vim split.

In `init.lua`:

```lua
-- https://github.com/croaky/laptop/blob/main/vim/init.lua

-- Helper functions
local function filetype_autocmd(ft, callback)
	vim.api.nvim_create_autocmd("FileType", {
		pattern = ft,
		callback = callback,
	})
end

local function run_file(key, cmd_template, split_cmd)
	local cmd = cmd_template:gsub("%%", vim.fn.expand("%:p"))
	buf_map(0, "n", key, function()
		vim.cmd(split_cmd)
		vim.cmd("terminal " .. cmd)
	end)
end

-- SQL
filetype_autocmd("sql", function()
	run_file("<Leader>r", "psql -d $(cat .db) -f % | less", "split")
end)
```

The `.db` file in the project contains only the local database name:

```
example_dev
```

See `man psql` for more detail on the `-d` and `-f` flags.
