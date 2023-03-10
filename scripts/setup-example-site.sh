#!/bin/bash

#
# This script sets up a local example site for testing,
# with Artalk integrated.
# Build frontend and backend before running this script.
# After running this script, run:
#   ./bin/artalk server -c ./local/local.yml
# to start the comment server.
# And open http://localhost:1313/ in your browser to view the example site.
#

# exit if any command fails
set -e

# check if hugo is installed
if ! [ -x "$(command -v hugo)" ]; then
  echo 'Error: hugo is not installed.' >&2
  exit 1
fi

# check if ./public/dist/Artalk.css exists
if ! [ -f "./public/dist/Artalk.css" ]; then
  echo 'Error: ./public/dist/Artalk.css does not exist.' >&2
  exit 1
fi

mkdir -p ./local

# copy internal/query/testdata/example_site_conf.yml to ./local/local.yml if it does not exist
if ! [ -f "./local/artalk.yml" ]; then
  echo "Copying internal/query/testdata/example_site_conf.yml to ./local/artalk.yml"
  cp ./internal/query/testdata/example_site_conf.yml ./local/local.yml
fi

# clean up local/example_site if it exists
if [ -d "./local/example_site" ]; then
  echo "Cleaning up local/example_site"
  rm -rfd ./local/example_site
fi

echo "Setting up local/example_site"
# create a new example site
hugo new site ./local/example_site

# copy the theme to the example site
theme_repo="https://github.com/ph-ph/chalk"
theme_dir="./local/example_site/themes/chalk"
echo "Cloning theme from ${theme_repo} to ${theme_dir}"
git clone ${theme_repo} ${theme_dir}


# create comment partial
echo "Creating comment partial in local/example_site"
mkdir -p ./local/example_site/layouts/partials/comments
cat << EOF > ./local/example_site/layouts/partials/comments/comments.html
<link href="/lib/artalk/Artalk.css" rel="stylesheet">
<script src="/lib/artalk/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>

<script>
  Artalk.init({
    el:        '#Comments',
    pageKey:   '{{ .Permalink }}',
    pageTitle: '{{ .Title }}',
    server:    '{{ $.Site.Params.artalk.server }}',
    site:      '{{ $.Site.Params.artalk.site }}',
  })
</script>
EOF

# patch ${single_template}, insert code before last </article> tag
single_template=local/example_site/themes/chalk/layouts/_default/single.html
echo "Patching ${single_template}"
sed -i '/<\/article>/i {{ partial "comments/comments.html" . }}' ./${single_template}

# update the example site config
cat << EOF > ./local/example_site/config.toml
baseURL = 'http://example.org/'
languageCode = 'en-us'
title = 'My New Hugo Site'
theme = "chalk"
[params.chalk]
  # chalk theme parameters
  about_enabled = true
  scrollappear_enabled = true
  theme_toggle = true
  rss_enabled = true
  blog_theme = 'light'
  local_fonts = false
[params.social]
  twitter='example'
  github='example'
[params.artalk]
  server='http://localhost:23366'
  site='Default Site'
EOF

# copy ./public/dist/Artalk.css/js to /lib/artalk/Artalk.css/js
echo "Copying ./public/dist/Artalk.css/js to /lib/artalk/Artalk.css/js"
mkdir -p ./local/example_site/static/lib/artalk
cp ./public/dist/Artalk.css ./local/example_site/static/lib/artalk/Artalk.css
cp ./public/dist/Artalk.js ./local/example_site/static/lib/artalk/Artalk.js

# create a post in the example site
echo "Creating a post in local/example_site"
mkdir -p ./local/example_site/content/posts
cat << EOF > ./local/example_site/content/posts/my-first-post.md
---
title: "My First Post"
---

# Manibus positaeque agrestes

## Inmurmurat posset saxa

Lorem markdownum undas, sors matre fata nimiumque parte; sua carae medios. Per
ego ad turritaque pars. Privignae tenuit auget in saxo nebulis. Pende enim
Silvanusque idem favilla, tecta! Motu *modo*, Almo vertuntur sceptroque alto,
mitia.

## Fuit ardor centum materiam in inquit

Rarescit digitis, turba debuit auctor inplerant monstra, pararet. Detestatur
**saltem dubitet Aeneadae** pectora celsis doloris, est tutus et circumdat agit
poposcit *precor* et.

## Bos avidus flammas excusare corpore

Nova Rutulos, Clytien, conspicuus poterit. Ne hac cladem fatendo postquam
saevitiae Dianae torvos *fluctus*, iuvenis et. Sensit adstitit separat effectum
saevarum fecit, me date imo, adicis nubigenas **ferumque**. Res illa signis!
Videri vagi, hac mori iras exuit similis modo et tempora.

    if (day_widget_saas) {
        footer.rate_dhcp(4, association_data, winSpoofingMca(real,
                esportsRepository, system));
    }
    if (unix_dma.serp(4 + trinitron, -3) ==
            koffice_function.defragmentQuadMatrix(vista_kernel,
            firmware_ntfs_tebibyte, -1 * bing)) {
        utility *= web;
        tunneling_halftone_primary(-2 + bare, -5, bmp + 5);
    }
    uddi += crop_error_capacity(bugCpc + boot, publishingSystemOutput(
            dvdClientOptical));
    if (2 == file_apple) {
        cutDevelopmentData.heuristic_bridge_desktop(
                softwareInfotainmentEdutainment - spam, tween,
                donationwareSecondary);
    } else {
        serverSdramRuntime = -2;
    }

## Nostrique miratur repetita suprema vocat sustulit

Dumosaque adflabat mixtaeque communicat est, Aeolus pinum. Me fecit.

    domain_leak_namespace.visualMonitor(header, gps_archie(service, 4,
            basic_networking_unfriend));
    cisc_software_cycle.fiosHertzAndroid += plug - 5 + storage(
            firewall_algorithm_troubleshooting);
    eup_botnet -= 5;
    computerSampleNull += uml * hard_clone_cd;

Limine habebat frigore, venenis tenebris nota. Fors erubui victa cortice
nullamque iniqui domos viroque, *inquit* ingentia, *crudelia*. Foedera qua!
*Ipsa dum perque*, mea orbi; aversus loca rutilos dona dixerat Bacchus saltus.
EOF


# run the example site
echo "Running local/example_site"
pushd ./local/example_site
hugo server
popd
