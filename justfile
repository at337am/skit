PROJECT_DIRS := "aout byebye cutimg dirhash fmp4 furl gostart img2vid jpg2png md2pg repaudio tmrn vcat vid2img vid2mp4 xfixer xnet"

default:
    @just --list

install-all:
    @echo "🚀 开始安装所有 skit 脚本..."
    @for dir in {{PROJECT_DIRS}}; do \
        echo ">>>>> 处理目录: $dir <<<<<"; \
        if [ -d "$dir" ]; then \
            ( \
                echo "  -> 切换到目录 '$dir'"; \
                cd "$dir" && \
                go install . && \
                echo "  ✅ '$dir' Successfully installed"; \
            ) || echo "  ❌ 错误：从 '$dir' 安装失败"; \
        else \
            echo "  💡 警告：目录 '$dir' 未找到。跳过。"; \
        fi; \
    done
    @echo "🎉 install-all 已完成"

