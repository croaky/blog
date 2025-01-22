# vim / format sql on save

In `laptop.sh`:

```bash
# https://github.com/croaky/laptop/blob/main/laptop.sh

brew install pgformatter
```

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

local function format_on_save(cmd_template)
	vim.api.nvim_create_autocmd("BufWritePre", {
		buffer = 0,
		callback = function()
			local cmd = cmd_template:gsub("%%", vim.fn.expand("%:p"))
			local buf = vim.api.nvim_get_current_buf()

			vim.fn.jobstart(cmd, {
				stdout_buffered = true,
				on_stdout = function(_, data)
					if not data then
						return
					end

					local fmt = table.concat(data, "\n")
					if #fmt == 0 then
						return
					end

					local pos = vim.api.nvim_win_get_cursor(0)

					local lines = vim.split(fmt, "\n")
					if lines[#lines] == "" then
						table.remove(lines, #lines)
					end

					vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)

					local line_count = vim.api.nvim_buf_line_count(buf)
					if pos[1] > line_count then
						pos[1] = line_count
					end

					vim.api.nvim_win_set_cursor(0, pos)

					vim.api.nvim_buf_call(buf, function()
						vim.cmd("noautocmd write")
					end)
				end,
				on_stderr = function(_, data)
					if data then
						print(table.concat(data, "\n"))
					end
				end,
			})
		end,
	})
end

-- SQL
filetype_autocmd("sql", function()
	format_on_save("pg_format --function-case 1 --keyword-case 2 --spaces 2 --no-extra-line %")
end)
```
