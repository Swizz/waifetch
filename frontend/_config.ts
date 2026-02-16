import site from "sprinkle/config.ts"

site.add("main.ts")

site.add("main.css")

site.copy("gh:loganmarchione/homelab-svg-assets@0.4.17/assets/*.svg", "loganmarchione/homelab-svg-assets/")

site.copy("wailsjs")

export default site
