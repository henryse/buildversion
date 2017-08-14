# Description

This program is used to generate version data we use for our services.

It accepts the following parameters:

* image - (required) Docker image for this application.
* imageid - (required) Docker tagged image id.
* output - (optional) File to output date to, the default is build_version.json

This will generate a build_version.json file we use for our **/version** calls on our services.

This is an [example](https://www.2829applegate.com/active):

    {
        "version": "c2692b51fc5f4e7bb7792b0764bdff48",
        "build_time": "2017-07-30_09:40:52",
        "image": "blog"
        "image_id": "17",
        "versions": {
            "debian": "52f04f4cf4eb4dc091ae5c2efceb7798",
            "nginx": "6b407ff8b85c4926ba56a780a18045ee",
            "proxy": "f0be97025af54421bf8074289061b104",
            "node": "f6eb68798add47b3b763fe44fa6b23fb",
            "ghost": "c2692b51fc5f4e7bb7792b0764bdff48"
        }
    }

This allows to to trace what build/version we are using for any given image via the URL.