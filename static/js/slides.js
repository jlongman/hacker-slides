function isPreview() {
    const preview = !!window.location.search.match(/preview/gi);
    console.log("Preview " + preview);
    return preview;
}


function initializeReveal() {
    // Full list of configuration options available at:
    // https://github.com/hakimel/reveal.js#configuration

    let deck = new Reveal({
        controls: true,
        progress: true,
        history: true,
        center: true,
        transition: 'slide', // none/fade/slide/convex/concave/zoom
        transitionSpeed: isPreview() ? 'fast' : 'default',
        embedded: isPreview() ? true : false,
        plugins: [
            RevealMarkdown,
            RevealHighlight,
            RevealSearch,
            RevealNotes,
            RevealMath,

            RevealAnything,
            RevealAudioSlideshow,
            RevealChalkboard,
            RevealChart,
            RevealCustomControls,
            RevealEmbedTweet,
            RevealFullscreen,

        ]

    })
    deck.initialize();

    themesCtrl();
    return deck;
}

function slide(slide) {
    deck.slide(slide);
}

function highlightAnyCodeBlocks() {
    $(document).ready(function () {
        $('pre code').each(function (i, block) {
            hljs.highlightBlock(block);
        });
    });
}

function insertMarkdownReference() {
    var markdownReference = $('<section/>', {
        'data-markdown': "/slides.md",
        'data-separator': "^---\n",
        'data-separator-vertical': "^--\n",
        'data-separator-notes': "^Note:",
        'data-charset': "utf-8"
    });

    $('.slides').html(markdownReference);
}

function scrollToCurrentSlide() {
    var i = deck.getIndices();
    deck.slide(i.h, i.v, i.f);
}

function reloadMarkdown() {
    insertMarkdownReference();
    deck.initialize();
    highlightAnyCodeBlocks();
    scrollToCurrentSlide();
}

function externalLinksInNewWindow() {
    $(document.links).filter(function () {
        return this.hostname != window.location.hostname;
    }).attr('target', '_blank');
}

insertMarkdownReference();
let deck = initializeReveal();

function themesCtrl() {
    var defaultTheme = "black.css",
        currentTheme = localStorage.getItem('theme?') ||
            defaultTheme;

    function setTheme(theme) {
        cssEl = $("#theme");
        cssEl.attr("href", "/static/reveal.js/css/theme/" + theme);
        localStorage.setItem('theme?', theme);
    }

    setTheme(currentTheme);

    if (!isPreview()) {
        return
    }
    var availableThemes = [
        "black.css",
        "beige.css",
        "blood.css",
        "league.css",
        "moon.css",
        "night.css",
        "serif.css",
        "simple.css",
        "sky.css",
        "solarized.css",
        "white.css",
    ];
    themeEl = $("#themes");
    availableThemes.forEach(function (theme) {
        elem = $("<option value=" + theme + ">" + theme + "</option>");
        themeEl.append(elem);
    })
    themeEl.val(currentTheme);
    themeEl.change(function () {
        val = themeEl.val()
        setTheme(val);
    });
    themeEl.attr("hidden", false);
}

// Monkey patch Reveal so we can reload markdown through an
// inter window message (using the reveal rpc api)
// (yes, reveal has an rpc api!)
// see save.js
deck.reloadMarkdown = reloadMarkdown;
