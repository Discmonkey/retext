<template>
    <div class="uploader">
        <div v-show="$refs.upload && $refs.upload.dropActive" class="drop-active">
            <h3>Drop files to upload</h3>
        </div>
        <table class="table table-hover">
            <thead>
            <tr>
                <th>#</th>
                <th>Thumb</th>
                <th>Name</th>
                <th>Size</th>
                <th>Speed</th>
                <th>Status</th>
                <th>Action</th>
            </tr>
            </thead>
            <tbody>
            <tr v-if="!files.length">
                <td colspan="7">
                    <div class="text-center p-5">
                        <h4>Drop files anywhere to upload<br/>or</h4>
                        <label :for="name" class="btn btn-lg btn-primary">Select Files</label>
                    </div>
                </td>
            </tr>
            <tr v-for="(file, index) in files" :key="file.id">
                <td>{{index}}</td>
                <td>
                    <img v-if="file.thumb" :src="file.thumb" width="40" height="auto" />
                    <span v-else>No Image</span>
                </td>
                <td>
                    <div class="filename">
                        {{file.name}}
                    </div>
                    <div class="progress" v-if="file.active || file.progress !== '0.00'">
                        <div :class="{'progress-bar': true, 'progress-bar-striped': true, 'bg-danger': file.error, 'progress-bar-animated': file.active}" role="progressbar" :style="{width: file.progress + '%'}">{{file.progress}}%</div>
                    </div>
                </td>
                <td>{{file.size }}</td>
                <td>{{file.speed}}</td>

                <td v-if="file.error">{{file.error}}</td>
                <td v-else-if="file.success">success</td>
                <td v-else-if="file.active">active</td>
                <td v-else></td>
                <td>
                    <div class="btn-group">
                        <button class="btn btn-secondary btn-sm dropdown-toggle" type="button">
                            Action
                        </button>
                        <div class="dropdown-menu">
                            <a :class="{'dropdown-item': true, disabled: file.active || file.success || file.error === 'compressing'}" href="#" @click.prevent="file.active || file.success || file.error === 'compressing' ? false :  onEditFileShow(file)">Edit</a>
                            <a :class="{'dropdown-item': true, disabled: !file.active}" href="#" @click.prevent="file.active ? $refs.upload.update(file, {error: 'cancel'}) : false">Cancel</a>

                            <a class="dropdown-item" href="#" v-if="file.active" @click.prevent="$refs.upload.update(file, {active: false})">Abort</a>
                            <a class="dropdown-item" href="#" v-else-if="file.error && file.error !== 'compressing' && $refs.upload.features.html5" @click.prevent="$refs.upload.update(file, {active: true, error: '', progress: '0.00'})">Retry upload</a>
                            <a :class="{'dropdown-item': true, disabled: file.success || file.error === 'compressing'}" href="#" v-else @click.prevent="file.success || file.error === 'compressing' ? false : $refs.upload.update(file, {active: true})">Upload</a>

                            <div class="dropdown-divider"></div>
                            <a class="dropdown-item" href="#" @click.prevent="$refs.upload.remove(file)">Remove</a>
                        </div>
                    </div>
                </td>
            </tr>
            </tbody>
        </table>

        <div class="btn-group">
            <file-upload
                    class="btn btn-primary dropdown-toggle"
                    :post-action="postAction"
                    :put-action="putAction"
                    :extensions="extensions"
                    :accept="accept"
                    :multiple="multiple"
                    :directory="directory"
                    :size="size || 0"
                    :thread="thread < 1 ? 1 : (thread > 5 ? 5 : thread)"
                    :headers="headers"
                    :data="data"
                    :drop="drop"
                    :drop-directory="dropDirectory"
                    :add-index="addIndex"
                    v-model="files"
                    @input-filter="inputFilter"
                    @input-file="inputFile"
                    ref="upload">
                <i class="fa fa-plus"></i>
                Select

            </file-upload>
            <div class="dropdown-menu">
                <label class="dropdown-item" :for="name">Add files</label>
                <a class="dropdown-item" href="#" @click="onAddFolder">Add folder</a>
                <a class="dropdown-item" href="#" @click.prevent="addData.show = true">Add data</a>
            </div>
            <button type="button" class="btn btn-success" v-if="!$refs.upload || !$refs.upload.active" @click.prevent="$refs.upload.active = true">
                <i class="fa fa-arrow-up" aria-hidden="true"></i>
                Start Upload
            </button>
            <button type="button" class="btn btn-danger"  v-else @click.prevent="$refs.upload.active = false">
                <i class="fa fa-stop" aria-hidden="true"></i>
                Stop Upload
        </button>
        </div>
    </div>
</template>

<script>
    import FileUpload from 'vue-upload-component'
    export default {
        components: {FileUpload},
        name: "Uploader",
        data() {
            return {
                files: [],
                accept: '*',
                extensions: 'txt,pdf',
                // extensions: ['gif', 'jpg', 'jpeg','png', 'webp'],
                // extensions: /\.(gif|jpe?g|png|webp)$/i,
                minSize: 1024,
                size: 1024 * 1024 * 10,
                multiple: true,
                directory: false,
                drop: true,
                dropDirectory: true,
                addIndex: false,
                thread: 3,
                name: 'file',
                postAction: '/file/upload',
                autoCompress: 1024 * 1024,
                uploadAuto: false,
                isOption: false,
                data: {},
                headers:{},
                addData: {
                    show: false,
                    name: '',
                    type: '',
                    content: '',
                },
                editFile: {
                    show: false,
                    name: '',
                }
            }
        },

        watch: {
            'addData.show'(show) {
                if (show) {
                    this.addData.name = ''
                    this.addData.type = ''
                    this.addData.content = ''
                }
            },
        },
        methods: {
            // add, update, remove File Event
            inputFile(newFile, oldFile) {
                if (newFile && oldFile) {
                    // update
                    if (newFile.active && !oldFile.active) {
                        // beforeSend
                        // min size
                        if (newFile.size >= 0 && this.minSize > 0 && newFile.size < this.minSize) {
                            this.$refs.upload.update(newFile, { error: 'size' })
                        }
                    }
                    if (newFile.progress !== oldFile.progress) {
                        // progress
                    }
                    if (newFile.error && !oldFile.error) {
                        // error
                    }
                    if (newFile.success && !oldFile.success) {
                        // success
                    }
                }
                if (!newFile && oldFile) {
                    // remove
                    if (oldFile.success && oldFile.response.id) {
                        // $.ajax({
                        //   type: 'DELETE',
                        //   url: '/upload/delete?id=' + oldFile.response.id,
                        // })
                    }
                }
                // Automatically activate upload
                if (Boolean(newFile) !== Boolean(oldFile) || oldFile.error !== newFile.error) {
                    if (this.uploadAuto && !this.$refs.upload.active) {
                        this.$refs.upload.active = true
                    }
                }
            },
            inputFilter(newFile, oldFile, prevent) {
                console.log(newFile, oldFile, prevent);
            },
            alert(message) {
                alert(message)
            },
            onEditFileShow(file) {
                this.editFile = { ...file, show: true }
                this.$refs.upload.update(file, { error: 'edit' })
            },
            onEditorFile() {
                if (!this.$refs.upload.features.html5) {
                    this.alert('Your browser does not support')
                    this.editFile.show = false
                    return
                }
                let data = {
                    name: this.editFile.name,
                }
                if (this.editFile.cropper) {
                    let binStr = atob(this.editFile.cropper.getCroppedCanvas().toDataURL(this.editFile.type).split(',')[1])
                    let arr = new Uint8Array(binStr.length)
                    for (let i = 0; i < binStr.length; i++) {
                        arr[i] = binStr.charCodeAt(i)
                    }
                    data.file = new File([arr], data.name, { type: this.editFile.type })
                    data.size = data.file.size
                }
                this.$refs.upload.update(this.editFile.id, data)
                this.editFile.error = ''
                this.editFile.show = false
            },
            // add folader
            onAddFolder() {
                if (!this.$refs.upload.features.directory) {
                    this.alert('Your browser does not support')
                    return
                }
                let input = this.$refs.upload.$el.querySelector('input')
                input.directory = true
                input.webkitdirectory = true
                this.directory = true
                input.onclick = null
                input.click()
                input.onclick = () => {
                    this.directory = false
                    input.directory = false
                    input.webkitdirectory = false
                }
            },
            onAddData() {
                this.addData.show = false
                if (!this.$refs.upload.features.html5) {
                    this.alert('Your browser does not support')
                    return
                }
                let file = new window.File([this.addData.content], this.addData.name, {
                    type: this.addData.type,
                })
                this.$refs.upload.add(file)
            }
        }
    }
</script>

<style>
    .uploader .btn-group .dropdown-menu {
        display: block;
        visibility: hidden;
        transition: all .2s
    }
    .uploader .btn-group:hover > .dropdown-menu {
        visibility: visible;
    }
    .uploader label.dropdown-item {
        margin-bottom: 0;
    }
    .uploader .btn-group .dropdown-toggle {
        margin-right: .6rem
    }
    .uploader .filename {
        margin-bottom: .3rem
    }

    .uploader .btn-is-option {
        margin-top: 0.25rem;
    }
    .uploader .example-foorer {
        padding: .5rem 0;
        border-top: 1px solid #e9ecef;
        border-bottom: 1px solid #e9ecef;
    }
    .uploader .edit-image img {
        max-width: 100%;
    }
    .uploader .edit-image-tool {
        margin-top: .6rem;
    }
    .uploader .edit-image-tool .btn-group{
        margin-right: .6rem;
    }
    .uploader .footer-status {
        padding-top: .4rem;
    }
    .uploader .drop-active {
        top: 0;
        bottom: 0;
        right: 0;
        left: 0;
        position: fixed;
        z-index: 9999;
        opacity: .6;
        text-align: center;
        background: #000;
    }
    .uploader .drop-active h3 {
        margin: -.5em 0 0;
        position: absolute;
        top: 50%;
        left: 0;
        right: 0;
        -webkit-transform: translateY(-50%);
        -ms-transform: translateY(-50%);
        transform: translateY(-50%);
        font-size: 40px;
        color: #fff;
        padding: 0;
    }
</style>
