PROJECT_DIRS := "dirhash gonew md2pg siho tmrn vid2mp4 xla"

default:
    @just --list

install-all:
    @for dir in {{PROJECT_DIRS}}; do \
        echo "Installing: $dir"; \
        if [ -d "$dir" ]; then \
            ( \
                cd "$dir" && \
                go install . && \
                echo "OK -> $dir"; \
            ) || echo "Erorr: '$dir' installation failed."; \
        else \
            echo "Error: '$dir' does not exist."; \
        fi; \
    done
    @echo "All skit CLI apps installed."
