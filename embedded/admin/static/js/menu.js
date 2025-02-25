document.addEventListener('DOMContentLoaded', function() {
    const menuToggle = document.getElementById('menu-toggle');
    const sidebar = document.getElementById('sidebar');
    
    function toggleMenu() {
        const isOpen = !sidebar.classList.contains('-translate-x-full');
        
        if (isOpen) {
            // Close menu
            sidebar.classList.add('-translate-x-full');
            document.body.classList.remove('overflow-hidden');
        } else {
            // Open menu
            sidebar.classList.remove('-translate-x-full');
            document.body.classList.add('overflow-hidden');
        }
    }

    menuToggle.addEventListener('click', toggleMenu);

    // Handle clicks outside menu on mobile
    document.addEventListener('click', function(event) {
        const isMobile = window.innerWidth < 768;
        const isOpen = !sidebar.classList.contains('-translate-x-full');
        const clickedOutside = !sidebar.contains(event.target) && !menuToggle.contains(event.target);
        
        if (isMobile && isOpen && clickedOutside) {
            toggleMenu();
        }
    });

    // Close menu when window is resized to desktop view
    window.addEventListener('resize', function() {
        if (window.innerWidth >= 768) { // md breakpoint
            sidebar.classList.add('-translate-x-full');
            document.body.classList.remove('overflow-hidden');
        }
    });
});
