import React from 'react'

class Comment extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			isEditing: false,
			body: props.body
		};

		this.handleEditChange = this.handleEditChange.bind(this);
	}


	onEdit() {
		this.setState({
			isEditing: true
		})
	}

	onCancelEdit() {
		this.setState({
			isEditing: false
		})
	}

	onUpdate() {
		// bubble up request
		this.props.onEdit(this.props.id, this.state.body);

		this.setState({
			isEditing: false
		})
	}

	handleEditChange(event) {
		this.setState({
			body:event.target.value
		})
	}


	editing() {
		return (
			<div>
				<input type="text" value={this.state.body} onChange={this.handleEditChange} />
				<br/>
				<div className="btn-group btn-group-xs">
					<button className="btn btn-warning" onClick={()=> this.onCancelEdit()}>cancel</button>
					<button className="btn btn-success" 
							onClick={()=> this.onUpdate()}>update</button>
				</div>
			</div>
			)
	}

	notEditing() {
		return (<div> 
				<p>{this.props.body}</p>
				<div className="btn-group btn-group-xs">
					<button className="btn btn-default" onClick={()=> this.onEdit()}>edit</button>
					<button className="btn btn-danger" onClick={() => this.props.onDelete(this.props.id)}>delete</button>
				</div>
				</div>
				);
	}

	bodyDisplay() {
		if (this.state.isEditing) {
			return this.editing();
		} else {
			return this.notEditing();
		}
	}

	render() {
		return (
			<div className="comment">
				<h1>{this.props.username} - {formateDate(this.props.modtime)}</h1>
				<br/>
				{this.bodyDisplay()}
			</div>
		);
	}
}


function formateDate(s) {
	return s.slice(0, 10);
}


export default Comment