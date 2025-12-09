PROJECT_DIRS := "dirhash gonew md2pg siho tmrn vid2mp4 xla"

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
    @echo ""
    @echo "Done."

