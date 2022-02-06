## OpenAPI plugin template

* #### Create a new OpenAPI based blink-plugin.

## READ FIRST

### Getting the Repo

#### Option 1: From Github
* Press `Use this template`
 
* Naming convention **blink-\<short-plugin-name>**

* Plugin repository should be public

* In the new repo go to
    - Settings→Webhooks and add `http://jenkins.blinkops.com:8443/generic-webhook-trigger/invoke/?token=<repo name>`

    - With 
Content Type: `application/json`

    -  This webhook will run master job automatically on every commit to master.

#### Option 2: Use Blink
* ##### Use this playbook:
#### https://github.com/blinkops/blink-playbooks/tree/main/git-create-repo-from-template

## Getting Started
* First in the `go.mod` make sure you are on the latest version of `blink-openapi-sdk`.
  
* Run `replace.bash <plugin-name>`

* Add `openapi.yaml`.

* Add actions in `mask.yaml` 
  * or optionally pass `""` in the PluginMetadata to list all the actions.

* build and register the plugin.

* `docker build . -f ./build/Dockerfile -t <plugin-name>:latest`
  
* generate readme with actions. --> `plugin.MakeMarkdown()` **will overide this readme**.


### Enjoy 😄
