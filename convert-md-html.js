// convert-md-html.js - Markdown to HTML conversion utilities

/**
 * Initialize the Markdown to HTML converter
 * Sets up event listeners and marked.js options
 */
function initMarkdownConverter() {
    // Elements
    const mdFile = document.getElementById('mdFile');
    const uploadBtn = document.getElementById('uploadBtn');
    const editBtn = document.getElementById('editBtn');
    const menuButton = document.getElementById('menuButton');
    const menuPanel = document.getElementById('menuPanel');
    const editorModal = document.getElementById('editorModal');
    const closeModal = document.getElementById('closeModal');
    const mdEditor = document.getElementById('mdEditor');
    const updateBtn = document.getElementById('updateBtn');
    const preview = document.getElementById('preview');
    
    // Set up marked options
    marked.setOptions({
        breaks: true,
        gfm: true,
        headerIds: true,
        headerPrefix: 'markdown-header-'
    });
    
    // Menu toggle
    menuButton.addEventListener('click', function() {
        this.classList.toggle('active');
        menuPanel.classList.toggle('active');
    });
    
    // Close menu when clicking outside
    document.addEventListener('click', function(e) {
        if (!menuButton.contains(e.target) && !menuPanel.contains(e.target)) {
            menuButton.classList.remove('active');
            menuPanel.classList.remove('active');
        }
    });
    
    // Handle file upload button
    uploadBtn.addEventListener('click', function() {
        mdFile.click();
    });
    
    // Handle Markdown file
    function handleFile(file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            const content = e.target.result;
            renderMarkdown(content);
            mdEditor.value = content;
            editBtn.classList.remove('hidden');

            // Close menu after file upload
            menuButton.classList.remove('active');
            menuPanel.classList.remove('active');
        };
        reader.readAsText(file);
    }

    // Handle file selection
    mdFile.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (!file) return;
        handleFile(file);
    });
    
    // Open modal editor
    editBtn.addEventListener('click', function() {
        editorModal.classList.add('active');
        
        // Close menu when opening editor
        menuButton.classList.remove('active');
        menuPanel.classList.remove('active');
    });
    
    // Close modal editor
    closeModal.addEventListener('click', function() {
        editorModal.classList.remove('active');
    });
    
    // Close modal when clicking outside
    editorModal.addEventListener('click', function(e) {
        if (e.target === editorModal) {
            editorModal.classList.remove('active');
        }
    });
    
    // Update preview from editor
    updateBtn.addEventListener('click', function() {
        const content = mdEditor.value;
        renderMarkdown(content);
        editorModal.classList.remove('active');
    });
    
    // Enable drag and drop
    document.body.addEventListener('dragover', function(e) {
        e.preventDefault();
        e.stopPropagation();
        document.body.style.backgroundColor = '#eaf2f8';
    });
    
    document.body.addEventListener('dragleave', function(e) {
        e.preventDefault();
        e.stopPropagation();
        document.body.style.backgroundColor = '#f6f8fa';
    });
    
    document.body.addEventListener('drop', function(e) {
        e.preventDefault();
        e.stopPropagation();
        document.body.style.backgroundColor = '#f6f8fa';
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            const file = files[0];
            if (file.name.endsWith('.md') || file.name.endsWith('.markdown') || file.name.endsWith('.txt')) {
                handleFile(file);
            }
        }
    });
}

/**
 * Render markdown content to HTML
 * @param {string} markdown - The markdown content to render
 */
function renderMarkdown(markdown) {
    const preview = document.getElementById('preview');
    const html = marked.parse(markdown);
    preview.innerHTML = '<article class="markdown-body">' + html + '</article>';
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', initMarkdownConverter);
