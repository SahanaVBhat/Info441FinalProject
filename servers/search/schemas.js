const mongoose = require('mongoose');
const Schema = require('mongoose').Schema;

const channel = new Schema({
	id: { type: Number, required: true, unique: true},
	name: { type: String, unique: true, default: 'general' },
    description: String,
    private: Boolean,
    members: {type:[{id:Number, email:String}]},
    createdAt: { type: Date, required: true},
    creator: {type:{id:Number, email:String}},
    editedAt: Date
});


const message = new Schema({
	id: {type:Number, required:true, unique:true},
    channelID: {type:Number, required:true},
    body: {type:String, required:true},
    createdAt: { type: Date, required: true},
    creator: {type:{id:Number, email:String}},
    editedAt: Date
});

module.exports = {
    channel,
    message
};