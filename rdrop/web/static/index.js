// 从服务器获取页面信息并触发显示。
async function fetchAndRenderPage() {
    // 状态显示的元素
    const promptElement = document.getElementById('prompt');

    try {
        const response = await fetch('/api/info');

        const info = await response.json();

        // throw new Error('这是一个模拟的测试错误！');

        if (!response.ok) {
            throw new Error(info.error || '服务器好像开小差了');
        }

        // 文件信息处理
        if (info.fileName && info.fileSize) {
            displayFileCard(info, 'shareable-file-card');

            // 获取当前页面的标题，并根据文件名称更新
            const currentTitle = document.title;
            document.title = `${info.fileName} - ${currentTitle}`;
        }

        // 描述信息处理
        if (info.description) {
            displayContent(info.description, 'description-message');
        }

        // 文本片段处理
        if (info.snippet) {
            displayContent(info.snippet, 'snippet-content');
        }

        // 最后隐藏提示元素
        promptElement.classList.add('is-hidden');
    } catch (error) {
        // 如果上面任何一步出错了，就在这里处理。
        console.error('错误:', error);
        promptElement.textContent = `啊哦, 出了点问题! ${error.message}`;
        promptElement.classList.add('error');
    }
}

// 显示共享文件卡片
function displayFileCard(info, elementId) {
    const targetDiv = document.getElementById(elementId);
    const template = document.getElementById('file-item-template');

    // 克隆模板内容。true 表示深度克隆（包含所有子节点）
    const fileItemClone = template.content.cloneNode(true);
    fileItemClone.querySelector('.file-name').textContent = info.fileName;
    fileItemClone.querySelector('.file-size').textContent = info.fileSize;

    // 指定下载文件的默认名称
    fileItemClone.querySelector('.file-item').setAttribute('download', info.fileName);

    // 将克隆并填充好的元素添加到页面
    targetDiv.appendChild(fileItemClone);
    targetDiv.classList.remove('is-hidden');
}

// 辅助函数, 给指定ID元素的文本内容, 并显示
function displayContent(content, elementId) {
    const element = document.getElementById(elementId);
    element.textContent = content;
    element.classList.remove('is-hidden');
}

// 从这里开始, 因为使用了 <script defer>, 可以确保 DOM 元素已加载完毕
fetchAndRenderPage();
