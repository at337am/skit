async function loadResources() {
  try {
    const response = await fetch("/api/resources");
    const data = await response.json();
    if (data.code === 0 && Array.isArray(data.data)) {
      document.getElementById("resources").innerHTML = generateTree(data.data);
    } else {
      document.getElementById("resources").innerHTML = `<div class="error">加载失败：${data.message || '数据格式不正确'}</div>`;
    }
  } catch (error) {
    document.getElementById("resources").innerHTML = `<div class="error">加载失败：${error.message}</div>`;
  }
}

function getFileIcon(filename) {
  const ext = filename.split('.').pop().toLowerCase();
  if (['mp4', 'webm', 'mov'].includes(ext)) return '🎥';
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)) return '🖼️';
  if (['mp3', 'wav', 'ogg'].includes(ext)) return '🎵';
  if (['pdf', 'doc', 'docx', 'txt'].includes(ext)) return '📄';
  return '📎';
}

function generateTree(resources) {
  // Check if resources is valid and an array
  if (!resources || !Array.isArray(resources)) {
    return "<div class='error'>无效的资源数据</div>";
  }

  let html = "<ul>";
  resources.forEach(item => {
    if (item.type === "directory" && Array.isArray(item.children)) {
      html += `<li class="directory">
          <div class="directory-header">
            <span class="directory-icon">📁</span>
            <span class="directory-name">${item.name || '未命名文件夹'}</span>
          </div>
          <div class="folder-content" style="display: none;">
            ${generateTree(item.children)}
          </div>
        </li>`;
    } else if (item.type === "file") {
      html += `<li class="file-item">
          <span class="file-icon">${getFileIcon(item.name || '')}</span>
          <a href="${item.url || '#'}" target="_blank" class="file-name">${item.name || '未命名文件'}</a>
        </li>`;
    }
  });
  html += "</ul>";
  return html;
}

// 添加目录折叠功能
document.addEventListener('click', (e) => {
  if (e.target.closest('.directory-header')) {
    const content = e.target.closest('.directory').querySelector('.folder-content');
    content.style.display = content.style.display === 'none' ? 'block' : 'none';
  }
});

// 检查API响应数据
function logResponse(data) {
  console.log("API Response:", data);
  if (!data || !data.data || !Array.isArray(data.data)) {
    console.error("数据格式错误：", data);
  }
  return data;
}

// 修改加载过程，添加日志
async function loadResourcesWithLogging() {
  try {
    console.log("开始加载资源...");
    const response = await fetch("/api/resources");
    console.log("API响应状态:", response.status);
    const data = await response.json();
    console.log("API响应数据:", data);

    if (data.code === 0 && Array.isArray(data.data)) {
      document.getElementById("resources").innerHTML = generateTree(data.data);
    } else {
      console.error("数据格式错误:", data);
      document.getElementById("resources").innerHTML = `<div class="error">加载失败：${data.message || '数据格式不正确'}</div>`;
    }
  } catch (error) {
    console.error("加载出错:", error);
    document.getElementById("resources").innerHTML = `<div class="error">加载失败：${error.message}</div>`;
  }
}

window.onload = loadResourcesWithLogging;
