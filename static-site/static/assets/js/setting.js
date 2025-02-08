document.addEventListener('DOMContentLoaded', () => {
    // load the setting from the local
    const savedSize = localStorage.getItem('fontSize');
    
    if (savedSize) {
        document.documentElement.style.fontSize = `${savedSize}px`;
    }
});