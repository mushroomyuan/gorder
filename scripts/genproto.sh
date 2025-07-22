#!/usr/bin/env bash

#  ！！！上面这个东西不是注释，叫做 "Shebang"（沙邦） 行，是操作系统识别脚本执行方式的特殊语法，必须写在第一行，放错位置就没用了。
#  类似的写法包括：
#  #!/usr/bin/env python3 用env动态查找 Python 3 来运行脚本
#  #!/usr/bin/env ruby 用env动态查找 Ruby 来运行脚本
#  #!/usr/bin/env bash 用env动态查找 Bash 来运行脚本
#  #!/usr/bin/env sh 用env动态查找 Shell 来运行脚本
#  #!/usr/bin/env php 用env动态查找 PHP 来运行脚本
#  #!/bin/bash 用绝对路径查找 Bash 来运行脚本

set -euo pipefail

# 这句 set -euo pipefail 是 Bash 脚本中非常重要的“安全选项组合”，目的是“只要有错误就立刻退出”，避免隐性 bug。
# 它有以下几个作用：
# 1. -e：如果命令执行失败，立即退出脚本
# 2. -u：如果使用未定义的变量，立即退出脚本
# 3. -o pipefail：如果管道中的命令执行失败，立即退出脚本
# 这些参数顺序可以打乱
# 例如：
# set -e -u -o pipefail 
# set -e -o pipefail -u
# set -o pipefail -e -u
# 注意：-o pipefail 是 Bash 4.0 版本之后才有的，所以如果使用的是旧版本的 Bash，需要手动添加。


shopt -s globstar
# 开启 **（双星号）通配符功能，使 ** 能递归匹配子目录中的所有文件。
# shopt：是 Shell Option 的缩写，用于设置 Bash 的行为。
# -s：表示“set”（启用）选项。
# globstar：开启 ** 通配符匹配功能。


if ! [[ "$0" =~ scripts/genproto.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi
# 这段的意思是如果当前脚本的路径不包含 scripts/genopenapi.sh，就报错退出，提示你“必须从项目根目录执行这个脚本”。
# "$0" 当前脚本的文件路径。
# [[ ... ]]这是 Bash 的扩展条件判断语法（比 [...] 功能更强），可以用来做字符串比较、通配符匹配、正则匹配等。
# =~ 是正则匹配运算符。
# （0~255）是 退出状态码，0 通常表示成功，255 是特殊值，表示“严重错误”，通常用于表示脚本执行失败。


source ./scripts/lib.sh
# 把 lib.sh 里面的所有代码当作当前脚本的一部分执行，包括变量、函数、别的命令。

API_ROOT="./api"
# 变量名不需要大写，但通常习惯用大写表示环境变量或全局常量，方便区分。

function dirs {
    dirs=()
    while IFS= read -r dir; do
        dirs+=("$dir")
    done < <(find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u)
    echo "${dirs[@]}"  # 使用 dirs 而非 dir
}

# dirs=() 定义一个空数组

# IFS是 Bash 中的一个特殊变量，代表 Internal Field Separator（内部字段分隔符）。它决定了 shell 在进行“单词分割”时用什么字符作为分隔符。
# 举个例子，
# 这里是用空格做分隔。
# line="apple banana cherry"
# IFS=' ' read -r a b c <<< "$line"
# echo "$a"  # apple
# echo "$b"  # banana
# echo "$c"  # cherry
# 这里是用逗号分隔。
# line="a,b,c"
# IFS=',' read -r x y z <<< "$line"
# echo "$x"  # a
# echo "$y"  # b
# echo "$z"  # c
# IFS= read -r line 清空 IFS：禁止自动分词整行读入

# < <(...) 把命令输出“伪装”成输入
# while ...; do
#   ...
# done < <(command)
# 这种写法叫做 Process Substitution（进程替代）

# read -r 是 Bash 中的一个命令，用于读取一行输入，而其中的 -r 选项表示：不要将反斜杠 \ 当作转义字符处理。
# 几乎所有读取外部数据的 Shell 脚本里都建议加上 -r，防止特殊字符被错误处理，写得更健壮可靠。

# find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u
# 就是所有 .proto 文件所在目录的最后一级目录名，去重后按字母排序。
# find . -type f -name "*.proto"：在当前目录（.）下查找所有 以 .proto 结尾的文件
# -type f：仅匹配文件
# -exec dirname {} \;
# 表示对找到的每个 .proto 文件，执行 dirname 命令（获取文件所在目录），然后通过 \; 结束。
# xargs -n1 basename | sort -u
# 对上面每行输出（目录路径），提取其“最后一级目录名”：
# xargs 是用于将标准输入（stdin）转成命令行参数的工具。
# 例如：echo "a b c" | xargs echo hello 会输出 hello a hello b hello c
# -n1 表示每次只取一个参数，即每次只处理一行输入。
# basename 是一个命令，用于从路径中提取“文件名”或最后一级目录名”。
# sort -u 是排序并去重。
# 最终得到所有 .proto 文件所在目录的最后一级目录名，去重后按字母排序。
# 关于管道符，| 是管道符，表示将前一个命令的输出作为后一个命令的输入。可以理解为流式处理命令，管道缓冲区64kb

function pb_files {
  pb_files=$(find . -type f -name '*.proto')
  echo "${pb_files[@]}"
}
# ${变量名[@]} 是 数组展开写法，表示将数组中的每个元素逐个列出。
# 但问题是：这里 pb_files 其实是个普通字符串变量（不是数组），因为 $(...) 默认是单个字符串。所以这句 echo "${pb_files[@]}" 实际上等效于 echo "$pb_files"。

function gen_for_modules() {
  local go_out="internal/common/genproto"
  if [ -d "$go_out" ]; then
    log_warning "found existing $go_out,cleaning all files under it"
    run rm -rf $go_out
  fi

  for dir in $(dirs);do
    local service="${dir:0:${#dir}-2}"
    local pb_files="${service}.proto"
    if [ -d "$go_out/$dir" ]; then
      log_warning "found existing $go_out,cleaning all files under it"
      run rm -rf "$go_out/$dir/*"
    else
      run mkdir -p "$go_out/$dir"
    fi
    log_info "generating code for $service to $go_out/$dir"
    run protoc\
      -I="/usr/local/include/" \
      -I="${API_ROOT}" \
      "--go_out=${go_out}" --go_opt=paths=source_relative \
      --go-grpc_opt=require_unimplemented_servers=false \
      "--go-grpc_out=${go_out}" --go-grpc_opt=paths=source_relative \
      "${API_ROOT}/${dir}/$pb_files"
  done
  log_success "proto gen done!"
}
# 这个函数是用来 遍历所有包含 .proto 文件的目录，并用 protoc 命令生成 Go 语言的 gRPC 代码，输出到 internal/common/genproto 目录下。
# local 是局部变量声明关键字
# source_relative 是 protoc 的官方参数选项之一，专门用于控制生成的代码文件相对于 .proto 文件的路径结构。

echo "directories containing protos to be built:$(dirs)"
echo "found pb_files:$(pb_files)"
gen_for_modules