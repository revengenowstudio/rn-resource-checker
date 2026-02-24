module main

import os
import time
import log

const sleep_second_at_job_done = 5

fn do_hash_job(target_paths []string, suffixes []string) ! {
	start_time := time.now()

	mut target_dirs := target_paths.clone()
	mut white_list := suffixes.clone()

	// 默认值处理
	if target_dirs.len == 0 {
		target_dirs = ['./']
	}
	if white_list.len == 0 {
		white_list = ['mix', 'exe', 'dll', 'ext']
	}

	log.info('Target paths: ${target_dirs}')
	log.info('Suffix whitelist: ${white_list}')

	// 扫描文件 (V 的 os.walk 非常快)
	mut files := []string{}
	for dir in target_dirs {
		log.info(dir)
        os.walk_with_context(dir, &files, fn [white_list] (mut f []string, path string) {
        ext := os.file_ext(path).replace('.', '')
        if !os.is_dir(path) && ext in white_list {
            f << path
        }
    })
    }
	log.info('Found ${files.len} files. Computing hashes...')

	// 并发计算
	results := compute_hashes_concurrently(files, 8)

	// 读取 Versioncode
	version_content := os.read_file('Versioncode') or {
		log.warn('Load Version file failed: ${err}')
		''
	}

	// 排序并输出
	mut sorted_results := results.clone()
	sorted_results.sort(a.filename < b.filename)

	timestamp := time.now().format_ss().replace(' ', '-').replace(':', '')
	out_name := 'result.${timestamp}.txt'

	mut out_file := os.create(out_name)!
	out_file.writeln(version_content)!
	for r in sorted_results {
		if !r.err {
			out_file.writeln('${r.filename} , hash : ${r.hash}')!
		}
	}
	out_file.writeln('Total file numbers : ${sorted_results.len}')!
	out_file.close()

	duration := time.since(start_time)
	log.info('Total elapsed time: ${duration}')

	println('\nWaiting for ${sleep_second_at_job_done} seconds before exiting...')
	time.sleep(sleep_second_at_job_done * time.second)
}
