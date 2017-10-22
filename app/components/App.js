import styles from './../style.css';
import React from 'react'
import axios from 'axios'

import CommentsList from './CommentsList';


class CommentsContainer extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			comments: []
		}

		//console.log(JSON.stringify(this.state.comments));

		this.onAdd    = this.onAdd.bind(this);
		this.onDelete = this.onDelete.bind(this);
		this.onEdit   = this.onEdit.bind(this);
	}

	initComments() {
		axios.get('http://localhost:8000/comments/1')
		.then ( res => {
			console.log("data received" + JSON.stringify(res.data));
			this.setState({
				comments: this.state.comments.concat(
					res.data
				)
			});

			console.log(this.state.comments);
		})
		.catch ( err => {
			console.log("Error:" + JSON.stringify(err))
		});
	}

	componentDidMount() {
		this.initComments();
	}

	onAdd(comment) {
		this.setState({
			comments: this.state.comments.concat([
				comment
			])
		});
	}

	onEdit(id, body) {

		// update database
		let jsonString = `{"id":${id}, "body":"${body}"}`;

		axios.put(`http://localhost:8000/comment/${id}`, jsonString, {
			headers:{ 'Content-Type': 'application/json'}
		})
		.then (res => {
			console.log(res.data)

			// create copy
			let newComments = this.state.comments;

			// mutate copy
			newComments.forEach(function(item, i) { 
				if (item.id == res.data.id)
					newComments[i] = res.data
			})

			// set state
			this.setState({
				comments: newComments
			});

		})
		.catch(err => console.log("Error: " + JSON.stringify(err)))
	}

	onDelete(id) {

		// delete from database
		axios.delete(`http://localhost:8000/comment/${id}`)
		.then ( res => {
			console.log(res.data);
		})
		.catch (err => console.log("Error: " + JSON.stringify(err)))

		let newcomments = this.state.comments.filter( el => el['id'] != id ); 
		this.setState({
			comments: newcomments
		});
	}

	render() {
		return (
			<div>
			<AddComment onAdd={this.onAdd} />
			<CommentsList comments={this.state.comments} onDelete={this.onDelete} onEdit={this.onEdit} />
			</div>
		)
	}
}




class AddComment extends React.Component{
	constructor(props) {
		super(props);

		this.state = {
			username: '',
			body: ''
		}

		this.handleUsernameChange = this.handleUsernameChange.bind(this);
		this.handleBodyChange = this.handleBodyChange.bind(this);
	}

	handleUsernameChange(event) {
		this.setState({username: event.target.value});
	}

	handleBodyChange(event) {
		this.setState({body: event.target.value});
	}

	handleOnAdd() {
			
		/* Add comment to database */
		var jsonString = `{"username": "${this.state.username}", "body": "${this.state.body}"}`;

		console.log("json string " + jsonString);

		axios.post('http://localhost:8000/comment', jsonString, {
			headers: { 'Content-Type': 'application/json'}
		})
		.then(res => {
			console.log(res.data);
			this.props.onAdd(res.data);
		})
		.catch( err => console.log("error: " + JSON.stringify(err.data)));
		
		this.setState({body: ''});
	}



	render() {
		return (
			<div className="comments-list add-comment">
				<h2> Add Comment </h2>
					<label>
						Username: 
						<input type="text" value={this.state.username} 
						 onChange={this.handleUsernameChange} />
					</label>
					<br/>
					<label>
						Comment:  
						<input type="text" value={this.state.body} 
						 onChange={this.handleBodyChange} />
					</label>
					<br/>
					<button className="btn btn-primary" onClick={() => this.handleOnAdd()}>Add Comment</button>
			</div>
		);
	}
}



class App extends React.Component {
	render() {
		return <CommentsContainer />
	}
}

export default App