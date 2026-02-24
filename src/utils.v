module main

import os
import crypto.sha256
import log
// import time

struct FileResult {
	filename string
	hash     string
	err      bool
}

// 拆分参数字符串
fn split_arg_string(input string) []string {
	mut sep := ' '
	if input.contains(',') {
		sep = ','
	} else if input.contains(';') {
		sep = ';'
	}
	return input.split(sep).map(it.trim_space())
}

// 计算单个文件的 SHA256
fn compute_sha256(path string) !string {
	mut file := os.open(path)!
	defer { file.close() }
	
	mut hash_obj := sha256.new()
	// V 的 io 读写很高效
	mut buffer := []u8{len: 131072}
	for {
		n := file.read(mut buffer) or { break }
		if n == 0 { break }
		hash_obj.write(buffer[..n]) or { break }
	}
	return hash_obj.sum([]u8{}).hex()
}

// // 并发计算 SHA256
// fn compute_hashes_concurrently(files []string, num_workers int) []FileResult {
// 	mut results := []FileResult{}
// 	ch := chan FileResult{cap: files.len}

// 	for path in files {
// 		spawn fn (path string, c chan FileResult) {
// 			hash := compute_sha256(path) or {
// 				c <- FileResult{filename: path, hash: '', err: true}
// 				return
// 			}
// 			// log.info('File ${path} : ${hash}')
// 			c <- FileResult{filename: path, hash: hash, err: false}
// 		}(path, ch)
// 	}

// 	for _ in 0 .. files.len {
// 		res := <-ch
// 		log.info('File ${res.filename} : ${res.hash}')
// 		results << res
// 	}
// 	return results
// }

fn compute_hashes_concurrently(files []string, num_workers int) []FileResult {
	mut results := []FileResult{}
	ch_task := chan string{cap: files.len}
	ch_res := chan FileResult{cap: files.len}

	// 1. 只开启固定数量的线程（比如 8 个）
	for _ in 0 .. num_workers {
		spawn fn (tasks chan string, res chan FileResult) {
			for {
				path := <-tasks or { break }
				h := compute_sha256(path) or {
					res <- FileResult{filename: path, hash: '', err: true}
					continue
				}
				res <- FileResult{filename: path, hash: h, err: false}
			}
		}(ch_task, ch_res)
	}

	// 2. 塞入任务
	for path in files { ch_task <- path }
	ch_task.close()

	// 3. 收集结果
	for _ in 0 .. files.len {
		res := <-ch_res
		results << res
		log.info('File ${res.filename} : ${res.hash}')
	}
	return results
}