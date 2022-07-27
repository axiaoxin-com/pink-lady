// https://github.com/matthornsby/priority-navigation
(function ($) {
  $.fn.prioritize = function (options) {
    $.fn.prioritize.defaults = {
      more: "More&hellip;",
      less: "Less&hellip;",
      menu: "Menu",
    };

    // Extend our default options with those provided.
    // Note that the first argument to extend is an empty
    // object – this is to keep from overriding our "defaults" object.
    var opts = $.extend({}, $.fn.prioritize.defaults, options);

    return this.each(function () {
      var element = this;

      $(element).addClass("prioritized");

      $('li[data-priority="0"]', element).addClass("demoted");

      var children = $(element).children(
        ':not(.demoted):not([data-priority="more"]):not([data-priority="less"]):not([data-priority="0"])'
      ).length;

      checkWidth(element);

      $(window).resize(function () {
        var children = $(element).children(
          ':not([data-priority="more"]):not([data-priority="less"]):not([data-priority="0"])'
        ).length;

        if (!$(element).hasClass("opened")) {
          $(element).removeClass("truncated");

          $(
            'li[data-priority="more"], li[data-priority="less"]',
            element
          ).remove();

          $('li:not([data-priority="0"])', element).removeClass("demoted");

          checkWidth(element, children);
        }
      });
    });

    function checkWidth(element, children) {
      var t = 0;

      //calculate the width of the visible <li>s
      $(element)
        .children()
        .not(".demoted")
        .outerWidth(function (i, w) {
          t += w;
        });

      if ($(element).css("display").indexOf("table") > -1) {
        var wrapper = $(element).parent().outerWidth();
      } else {
        // var wrapper = $(element).outerWidth();
        var wrapper = $(element).parent().outerWidth();
      }

      if (wrapper < t) {
        if (!$('li[data-priority="more"]', element).length) {
          $(element).append(
            '<li data-priority="more"><a href="#">' +
              opts.more +
              '</a></li><li data-priority="less"><a href="#">' +
              opts.less +
              "</a></li>"
          );
          //console.log("no");
        }

        hideTheHeighest(element, options, children);

        moreOrLess(element);
      }
    }

    function moreOrLess(element, children) {
      $('li[data-priority="more"] a', element).on("click", function (event) {
        event.preventDefault();
        //console.log("click");
        $(this).parents("ul").addClass("opened");
      });

      $('li[data-priority="less"] a', element).on("click", function (event) {
        event.preventDefault();
        //console.log("click");
        $(this).parents("ul").removeClass("truncated opened");
        $(
          'li[data-priority="more"], li[data-priority="less"]',
          element
        ).remove();
        $('li:not([data-priority="0"])', element).removeClass("demoted");
        checkWidth(element, children);
      });
    }

    function hideTheHeighest(element, options) {
      //console.log(children);

      if (
        $(element).children(
          ':not(.demoted):not([data-priority="more"]):not([data-priority="less"])'
        ).length < 2
      ) {
        $('[data-priority="more"]', element)
          .addClass("menu")
          .children()
          .text(opts.menu);
      } else {
        $('[data-priority="more"]', element).removeClass("menu");
      }

      $(element).addClass("truncated");

      /*
			//this hides the leftmost instance of the highest visible data-priority
			var max = 0, index = 0;
			$('.truncated > *:not(.demoted):not([data-priority="more"])').each(function(i){
			   if(parseInt($(this).data('priority'), 10) > max){
			      max = parseInt($(this).data('priority'), 10);
			      index = i;
			   }

			}).eq(index).addClass("demoted");
			*/

      //hides all of the highest visible data-priority, which is better, but has some resize issues in chrome
      var highestVisible = 0;
      $("*:not(.demoted)", element).each(function () {
        if ($.isNumeric($(this).data("priority"))) {
          if (parseInt($(this).data("priority"), 10) > highestVisible) {
            highestVisible = parseInt($(this).data("priority"), 10);

            // console.log("highest: " + highestVisible);
          }
        }
      });
      $('[data-priority="' + highestVisible + '"]', element).addClass(
        "demoted"
      );

      checkWidth(element);
    }
  };
})(jQuery);
