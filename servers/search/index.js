const mongoose = require('mongoose');
const express = require('express');
const { channel, message } = require('./schemas');
const handlers = require('./handlers');
const amqp = require('amqplib/callback_api');

// you can have multiple databases similar to sql.
// this is the test database
const rabbAddr = process.env.RABBITADDR || "amqp://rabbit:5672"
const rabbQueueName = process.env.RABBITQUEUENAME || "queue"

const mongoEndpoint = "mongodb://customMongoContainer:27017/rabbit";
//const rabbAddr = process.env.RABBITADDR || "amqp://localhost:5672";
//const mongoEndpoint = "mongodb://localhost:27017/test";
const port = 80;
let rabbitChannel;
// set up mongoose schemas
const Channel = mongoose.model("Channel", channel);
const Messages = mongoose.model("MessagesModel", message);

// set up express
const app = express();
app.use(express.json());


const getRabbitChannel = () => {
	return rabbitChannel;
}

// A function to connect to the mongo endpoint, used for refreshing on disconnect.
const connect = () => {
	mongoose.connect(mongoEndpoint);
}

//default 'general' channel
const createdAt = new Date();
const defaultChannel = { 
	id: 1,
	name: "general", 
	description: "default channel", 
	private: false, 
	members: [], 
	createdAt: createdAt, 
	creator: {}, 
	editedAt: "" 
};

const query = new Channel(defaultChannel);
query.save((err, newChannel) => {
	if (err) {
		console.log("Unable to create default channel 'general'.");
		return;
	}
	console.log("default channel created");
});

app.all('*',function(req,res,next)
{
    if (!req.get('Origin')) return next();

    res.set('Access-Control-Allow-Origin','http://myapp.com');
    res.set('Access-Control-Allow-Methods','GET,POST,DELETE,PATCH');
    res.set('Access-Control-Allow-Headers','X-Requested-With,Content-Type');

    if ('OPTIONS' == req.method) return res.send(200);

    next();
});

// for all channels
app.get("/v1/channels", async (req, res) => {
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
    try {
        const channels = await Channel.find([{"private":false}]);
        res.set("Content-Type", "application/json");
        res.json(channels);
    } catch (e) {
        res.status(500).send("There was an issue getting channels");
    }
});

app.post("/v1/channels", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	const { name, description, private, members, creator } = req.body;
	if (!name) {
		res.status(400).send("Must provide name of the channel");
		return;
	}
	//get number of documents in channel
	const Lastid = await Channel.countDocuments({});
	const id = Lastid+1;
	//get user
	var usr = JSON.parse(XUser);
 	const createdAt = new Date();

	const channel = { 
		id: id,
		name: name, 
		description: description, 
		private: private, 
		members: members, 
		createdAt: createdAt, 
		creator: usr
	};

	const query = new Channel(channel);
	query.save((err, newChannel) => {
		if (err) {
			console.log(err);
			res.status(500).send("Unable to create channel");
			return;
        }
        
        res.set("Content-Type", "application/json");
		res.status(201).json(newChannel);

		let userIDs = [];
		if (private) {
			//get userIDlist
			for (let m in members){
				userIDs.push(m.id);
			}
		}
					  
		let data = {
			type: "channel-new",
			channel: newChannel,
			userIDs: userIDs
		};

		rabbitChannel.sendToQueue(rabbQueueName, Buffer.from(JSON.stringify(data)));
	});
});

// for a specific channel identified by {channelID}
app.get("/v1/channels/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get channel id 
	const channelId = req.params.id;
	const specificChannel = await Channel.find({ id: channelId });
	//get user
	var usr = JSON.parse(XUser);
	//if private and user not member
	if (specificChannel[0].private){
		if (!specificChannel[0].members.includes(usr.id)){
		res.status(403).send("Forbidden User");
		return;
		}
	}

    var messageList = await Messages.find({ channelID: channelId }).sort({ createdAt: -1 }).limit(100);
    // the most recent 100 messages in the specified channel before a specific message
    const beforeMessageID = req.query.before;
    if (beforeMessageID) {
        messageList = await Messages.find({ channelID: channelId , id: {$lt:beforeMessageID}}).sort({ createdAt: -1 }).limit(100);
    }
    res.set("Content-Type", "application/json");
	res.status(201).json(messageList);
	
});

app.post("/v1/channels/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get channel id 
	const channelId = req.params.id;
	const specificChannel = await Channel.find({ id: channelId });
	//get user
	var usr = JSON.parse(XUser);
	//if private and user not member
	if (specificChannel[0].private){
		if (!specificChannel[0].members.includes(usr.id)){
		res.status(403).send("Forbidden User");
		return;
		}
	}
	//get number of documents in Messages
	const Lastid = await Messages.countDocuments({});
	const id = Lastid+1;
	// get message
	const {body}  = req.body;
 	const createdAt = new Date();
	const message = { 
		id : id,
		channelID: channelId, 
		body: body, 
		createdAt: createdAt, 
		creator: usr
	};

	const query = new Messages(message);
	query.save((err, newMessage) => {
		if (err) {
			res.status(500).send("Unable to create message");
			return;
        }
        res.set("Content-Type", "application/json");
		res.status(201).json(newMessage);

		let userIDs = [];
		if (specificChannel[0].private) {
			//get userIDlist
			for (let m in specificChannel[0].members){
				userIDs.push(m.id);
			}
		}
					  
		let data = {
			type: "message-new",
			message: newMessage,
			userIDs: userIDs
		};

		rabbitChannel.sendToQueue(rabbQueueName, Buffer.from(JSON.stringify(data)));
	});
});

app.patch("/v1/channels/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get channel id 
	const channelId = req.params.id;
	const specificChannel = await Channel.find({ id: channelId });
	//get user
	var usr = JSON.parse(XUser);
	//if private and user not member
	if (specificChannel[0].private){
		if (!specificChannel[0].members.includes(usr.id)){
		res.status(403).send("Forbidden User");
		return;
		}
	}

	if(specificChannel[0].creator.id != usr.id){
		res.status(403).send("Non-creator cannot modify channel");
		return;
	}
	//get updates 
	const {name, description} = req.body;
	// update name
	if(name) {
		await Channel.where({ id: channelId }).update({ name: name });
	}
	// update description
	if(description) {
		await Channel.where({ id: channelId }).update({ description: description });
	}
	const updatedChannel = await Channel.find({ id: channelId });
	res.set("Content-Type", "application/json");
	res.json(updatedChannel);

	let userIDs = [];
	if (specificChannel[0].private) {
		//get userIDlist
		for (let m in specificChannel[0].members){
			userIDs.push(m.id);
		}
	}
					
	let data = {
		type: "channel-update",
		channel: updatedChannel,
		userIDs: userIDs
	};

	rabbitChannel.sendToQueue(rabbQueueName, Buffer.from(JSON.stringify(data)));
		
});

app.delete("/v1/channels/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get channel id 
	const channelId = req.params.id;
	const specificChannel = await Channel.find({ id: channelId });
	//get user
	var usr = JSON.parse(XUser).id;
	//if private and user not member
	if (specificChannel[0].private){
		if (!specificChannel[0].members.includes(usr.id)){
		res.status(403).send("Forbidden User");
		return;
		}
	}

	if(specificChannel[0].creator.id != usr.id){
		res.status(403).send("Non-creator cannot delete channel");
		return;
	}

    if (specificChannel[0].name != "general") {
        // delete channel
        await Channel.deleteOne({ id: channelId });
        // delete messages for channel
        await Messages.deleteMany({ channelID: channelId });
    }
	
    res.set("Content-Type", "text/plain");
	res.send("Successful deletion");


	let userIDs = [];
	if (specificChannel[0].private) {
		//get userIDlist
		for (let m in specificChannel[0].members){
			userIDs.push(m.id);
		}
	}
					
	let data = {
		type: "channel-delete",
		channelID: channelId,
		userIDs: userIDs
	};

	rabbitChannel.sendToQueue(rabbQueueName, Buffer.from(JSON.stringify(data)));

});

// for the members of a private channel identified by {channelID}
app.post("/v1/channels/:id/members", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	
	//get channel id 
	const channelId = req.params.id;
	const specificChannel = await Channel.find({ id: channelId });
	//get user
	var usr = JSON.parse(XUser);
	//if private and user not member
	if (specificChannel[0].private){
		if (!specificChannel[0].members.includes(usr.id)){
		res.status(403).send("Forbidden User");
		return;
		}
	}
	if(specificChannel[0].creator.id != usr.id){
		res.status(403).send("Non-creator cannot add members");
		return;
	}
	//get updates 
	const {id, email} = req.body;
	const newUser = {
		id: id,
		email: email
	};
	const membersList = specificChannel[0].members;
	membersList.push(newUser);
	if(membersList) {
		await Channel.where({ id: channelId }).update({ members: membersList });
	}

    res.set("Content-Type", "text/plain");
	res.status(201).send("Success!! User has been added as a member");
});

app.delete("/v1/channels/:id/members", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	
	//get channel id 
	const channelId = req.params.id;
	const specificChannel = await Channel.find({ id: channelId });
	//get user
	var usr = JSON.parse(XUser);
	//if private and user not member
	
	const memberlist = specificChannel[0].members;
	if (specificChannel[0].private){
		if (!specificChannel[0].members.includes(usr.id)){
		res.status(403).send("Forbidden User");
		return;
		}
	}
	if(specificChannel[0].creator.id != usr.id){
		res.status(403).send("Non-creator cannot remove members");
		return;
	}
	
	//get updates 
	const {id, email} = req.body;
	var index = -1;
	//  remove user ID from array 
	for (let m in memberlist){
		const member = memberlist[m].toString();
		if (member.includes(id) && member.includes(email)){
			index = m;
		}
	}
	//var index = specificChannel[0].members.indexOf(usr);
	if (index !== -1) {
		specificChannel[0].members.splice(index, 1);
	}
	const memberslistUpdated = specificChannel[0].members;
	// update members
	await Channel.where({ id: channelId }).updateOne({ members: memberslistUpdated });

    res.set("Content-Type", "text/plain");
	res.status(200).send("Successfully removed user from members");
});

// for a specific message identified by {messageID}
app.patch("/v1/messages/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get message id 
	const messageId = req.params.id;
	const specificMessage = await Messages.find({ id: messageId });
	//get user
	var usr = JSON.parse(XUser);
	//if user not creator of message
	if (specificMessage[0].creator.id != usr.id) {
		res.status(403).send("User cannot edit post they didn't post");
		return;
	}

 	const editedAt = new Date();
	const message = { 
		channelID: specificMessage[0].channelID, 
		body: req.body, 
		createdAt: specificMessage[0].createdAt, 
		creator: usr,
		editedAt: editedAt
	};

	await Messages.where({ id: messageId }).updateOne({message});
	const updatedMessage = await Messages.find({ id: messageId });
	res.set("Content-Type", "application/json");
	res.json(updatedMessage);


	const specificChannel = await Channel.find({ id: updatedMessage[0].channelId });

	let userIDs = [];
	if (specificChannel[0].private) {
		//get userIDlist
		for (let m in specificChannel[0].members){
			userIDs.push(m.id);
		}
	}
					
	let data = {
		type: "message-update",
		message: updatedMessage,
		userIDs: userIDs
	};

	rabbitChannel.sendToQueue(rabbQueueName, Buffer.from(JSON.stringify(data)));

});

app.delete("/v1/messages/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get message id 
	const messageId = req.params.id;
	const specificMessage = await Messages.find({ id: messageId });
	//get user
	var usr = JSON.parse(XUser);
	//if user not creator of message
	if (specificMessage[0].creator.id != usr.id) {
		res.status(403).send("User cannot delete messages they didn't post");
		return;
	}
	// delete message
    await Messages.deleteOne({ id: messageId });
    res.set("Content-Type", "text/plain");
	res.send("Successfullly deleted message");

	const specificChannel = await Channel.find({ id: updatedMessage[0].channelId });

	let userIDs = [];
	if (specificChannel[0].private) {
		//get userIDlist
		for (let m in specificChannel[0].members){
			userIDs.push(m.id);
		}
	}
					
	let data = {
		type: "message-delete",
		messageID: messageId,
		userIDs: userIDs
	};

	rabbitChannel.sendToQueue(rabbQueueName, Buffer.from(JSON.stringify(data)));
});

const RequestWrapper = (handler, SchemeAndDbForwarder) => {
	return (req, res) => {
		handler(req, res, SchemeAndDbForwarder);
	}
}

app.post("/v1/channels", RequestWrapper(handlers.postChannelsHandler, { Channel, getRabbitChannel }));
app.get("/v1/channels", RequestWrapper(handlers.getChannelsHandler, { Channel }));

connect();
mongoose.connection.on('error', console.error)
	.on('disconnected', connect)
	.once('open', main);

async function main() {
	amqp.connect(rabbAddr, (err, conn) => {
		if (err) {
			console.log("Failed to connect to rabbit instance");
			process.exit(1);
		}

		conn.createChannel((err, ch) => {
			if (err) {
				console.log("Error creating channel")
				process.exit(1);
			}

			ch.assertQueue("queue", { durable: true });
			rabbitChannel = ch;


			ch.consume("queue", (msg) => {
				console.log(msg.content.toString());
			}, {
				noAck: true
			})
		});

		app.listen(port, "", () => {
			console.log(`Server listening on port ${port}`);
		});

		
	});
}
