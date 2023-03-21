/**********************************************************************
░██████╗░█████╗░██████╗░██╗██████╗░████████╗░██████╗░░░░░░░░██╗░██████╗
██╔════╝██╔══██╗██╔══██╗██║██╔══██╗╚══██╔══╝██╔════╝░░░░░░░░██║██╔════╝
╚█████╗░██║░░╚═╝██████╔╝██║██████╔╝░░░██║░░░╚█████╗░░░░░░░░░██║╚█████╗░
░╚═══██╗██║░░██╗██╔══██╗██║██╔═══╝░░░░██║░░░░╚═══██╗░░░██╗░░██║░╚═══██╗
██████╔╝╚█████╔╝██║░░██║██║██║░░░░░░░░██║░░░██████╔╝██╗╚█████╔╝██████╔╝
╚═════╝░░╚════╝░╚═╝░░╚═╝╚═╝╚═╝░░░░░░░░╚═╝░░░╚═════╝░╚═╝░╚════╝░╚═════╝░
**********************************************************************/
$(function(){
    $('.content-section').each(function() {
        if ($(window).height() > $(this).offset().top + 200) {
            $(this).addClass('loaded');
        }
    });

    $('.header-list-li').click(function() {
        var index = jQuery(this).index()
        var link = $('.header-list-li').eq(index).children('a').attr("href");
        location.href = link;
    });

    $('.header-title-box').click(function() {
        location.href = "/";
    });

    $('.ext-link1').click(function() {
        // content
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