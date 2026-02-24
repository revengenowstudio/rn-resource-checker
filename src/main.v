module main

import os
import cli
import log

fn main() {
	// 初始化日志格式
	log.set_level(.info)

	mut app := cli.Command{
		name: 'rn-resource-checker'
		description: 'Input path list and suffix list to compute SHA256'
		execute: fn (cmd cli.Command) ! {
			mut paths := cmd.flags.get_strings('path') or { []string{} }
			mut suffixes := cmd.flags.get_strings('suffix') or { []string{} }
			
			// 如果用户是通过逗号分隔输入的单条字符串，处理一下
			mut final_paths := if paths.len == 1 { split_arg_string(paths[0]) } else { paths }
			mut final_suffixes := if suffixes.len == 1 { split_arg_string(suffixes[0]) } else { suffixes }

			do_hash_job(final_paths, final_suffixes)!
		}
	}

	app.add_flag(cli.Flag{
		flag: .string_array
		name: 'path'
		abbrev: 'p'
		description: 'Target directory paths'
	})

	app.add_flag(cli.Flag{
		flag: .string_array
		name: 'suffix'
		abbrev: 's'
		description: 'File suffixes to match'
	})

	app.setup()
	app.parse(os.args)
}