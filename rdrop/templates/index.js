/**
 * 主要任务：从服务器获取文件信息并触发显示。
 * "async" 表示这个函数里有需要“等待”的操作。
 * @param {HTMLElement} targetDiv - 要在哪个HTML元素里显示内容。
 */
async function fetchAndShowFile(targetDiv) {
    try {
        // 去服务器的 '/api/files' 地址拿数据，并等待结果。
        const response = await fetch('/api/files');

        // 如果服务器没成功返回数据，就抛出一个易懂的错误。
        if (!response.ok) {
            throw new Error('啊哦，服务器好像开小差了，获取文件信息失败！');
        }

        // 把返回的 JSON 数据“翻译”成 JS 对象。
        const fileInfo = await response.json();

        // 调用另一个函数，把文件信息显示在指定的盒子里。
        displayFileOnPage(fileInfo, targetDiv);

    } catch (error) {
        // 如果上面任何一步出错了，就在这里处理。
        console.error('出错了:', error); // 在控制台打印详细错误。
        targetDiv.innerHTML = `😢 ${error.message}`; // 在网页上显示友好的错误提示。
        targetDiv.className = 'error';
    }
}

/**
 * 辅助任务：把文件信息创建成网页元素并显示出来。
 * @param {object} file - 包含文件信息的对象，如 { Name: '...', Size: '...' }。
 * @param {HTMLElement} targetDiv - 要把创建好的元素放进哪个父容器。
 */
function displayFileOnPage(file, targetDiv) {
    // 1. 清空原有的 "Loading..." 文字。
    targetDiv.innerHTML = '';
    targetDiv.classList.remove('loading');

    // 1. 获取模板
    const template = document.getElementById('file-item-template');
    
    // 2. 克隆模板内容。true 表示深度克隆（包含所有子节点）
    const fileItemClone = template.content.cloneNode(true);

    // 3. 填充数据 (使用 querySelector 在克隆的片段中查找元素)
    // 同样使用 textContent 保证安全
    fileItemClone.querySelector('.file-name').textContent = file.Name;
    fileItemClone.querySelector('.file-size').textContent = file.Size;

    // 4. 将克隆并填充好的元素添加到页面
    targetDiv.appendChild(fileItemClone);
}

// 因为HTML里用了 <script defer>, 浏览器保证这个脚本执行时，
// 整个HTML文档已经加载完毕。所以我们不再需要 'DOMContentLoaded' 包裹。

// 1. 找到那个我们要操作的“墙”（ID为'resources'的div）。
const resourcesDiv = document.getElementById('resources');

// 2. 调用我们之前定义好的函数，开始执行任务！
// 把我们找到的 div 作为参数传进去。
fetchAndShowFile(resourcesDiv);
