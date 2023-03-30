/**********************************************************************
░██████╗░█████╗░██████╗░██╗██████╗░████████╗░██████╗░░░░░░░░██╗░██████╗
██╔════╝██╔══██╗██╔══██╗██║██╔══██╗╚══██╔══╝██╔════╝░░░░░░░░██║██╔════╝
╚█████╗░██║░░╚═╝██████╔╝██║██████╔╝░░░██║░░░╚█████╗░░░░░░░░░██║╚█████╗░
░╚═══██╗██║░░██╗██╔══██╗██║██╔═══╝░░░░██║░░░░╚═══██╗░░░██╗░░██║░╚═══██╗
██████╔╝╚█████╔╝██║░░██║██║██║░░░░░░░░██║░░░██████╔╝██╗╚█████╔╝██████╔╝
╚═════╝░░╚════╝░╚═╝░░╚═╝╚═╝╚═╝░░░░░░░░╚═╝░░░╚═════╝░╚═╝░╚════╝░╚═════╝░
**********************************************************************/
$(function(){
    // overlay
    function overlay_action(elem, action) {
        $('.overlay-button-close').on('click', function() {
            $('.ext-link-overlay').fadeOut();
            return false;
        });
        $('.overlay-screen').click(function() {
            $('.ext-link-overlay').fadeOut();
            return false;
        });
        $('.overlay-box').click(function(event){
            event.stopPropagation();
            return false;
        });
        $('.overlay-button-action').click(function() {
            $('.ext-link-overlay').fadeOut();
            action(elem);
            return true;
        });
    }

    $('.content-section').each(function() {
        if ($(window).height() > $(this).offset().top + 200) {
            $(this).addClass('loaded');
        }
    });

    // make header list link
    $('.header-list-li').click(function() {
        var index = jQuery(this).index()
        var link = $('.header-list-li').eq(index).children('a').attr("href");
        location.href = link;
    });

    $('.header-title-box').click(function() {
        location.href = "/";
    });

    // ext link overlay
    $('.ext-link').click(function() {
        $(".ext-link-overlay").fadeIn();
        action = function(elem) {
            window.open(elem, "_blank");
            return false;
        }
        result = overlay_action(this, action);
        return false;
    });

    // when scrolled
    $(window).scroll(function () {
        // header
        if ($(window).scrollTop() - $('main').offset().top > 50) {
            $('header').css("position", "fixed");
        } else {
            $('header').css("position", "static");
        }

        // display content
        $('.content-section').each(function() {
            var pos = $(this).offset().top;
            var scroll = $(window).scrollTop();
            var height = $(window).height();
            if (scroll > pos - height + 200) {
                $(this).addClass('loaded');
            }
        });
    });
});