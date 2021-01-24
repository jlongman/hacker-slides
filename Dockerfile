FROM hacker-slides-go:latest AS compiler
ENV VERSION=1.0.0

FROM node:12-slim as revealjs
# Grab reveal.js code
ENV REVEAL_VERSION=4.1.0 \
    REPO=https://codeload.github.com/hakimel/reveal.js
#    SHA=45dc8caeb1a1a81d74293552ea3a408bc463dc3e
RUN apt-get update \
    && apt-get install -y wget
RUN wget -O /tmp/reveal.js.tar.gz  https://codeload.github.com/hakimel/reveal.js/tar.gz/$REVEAL_VERSION && \
    tar -xzf /tmp/reveal.js.tar.gz -C / && \
    rm -f /tmp/reveal.js.tar.gz && \
    mv reveal.js-$REVEAL_VERSION /revealjs
# skip sha1 check
#    echo "$SHA /tmp/reveal.js.tar.gz" | sha1sum --check - && \

# Switch to the reveal.js directory.
WORKDIR /revealjs

# Prepare the reveal.js installation.
RUN npm install && npm install --global gulp-cli
RUN gulp build

# install plugins, must still be added to source
#RUN npm install reveal.js-plugins && cp -pr node_modules/reveal.js-plugins-plugins/* plugin

ENV REVEAL_PLUGINS_VERSION=4.0.2
RUN wget -O /tmp/reveal.js-plugins.tar.gz  https://codeload.github.com/rajgoel/reveal.js-plugins/tar.gz/$REVEAL_PLUGINS_VERSION && \
    tar -xzf /tmp/reveal.js-plugins.tar.gz --strip-components=1 -C plugin && \
    rm -f /tmp/reveal.js-plugins.tar.gz


# core
FROM alpine:3.8

WORKDIR /srv

ENV GIN_MODE=release
RUN mkdir slides
COPY --from=compiler /bin/app /bin/app
COPY static static
COPY --from=revealjs /revealjs/plugin/   static/reveal.js/plugin/
COPY --from=revealjs /revealjs/dist/*.css /revealjs/dist/theme/ static/reveal.js/css/
COPY --from=revealjs /revealjs/dist/*.js*  static/reveal.js/js/

COPY templates templates
COPY initial-slides.md initial-slides.md


CMD app $PORT
