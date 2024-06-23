import React, { useState, FormEvent } from 'react';

const CreateTag: React.FC = () => {
  const [name, setName] = useState<string>('');

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const response = await fetch('/tags', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name }),
    });
    if (response.ok) {
      alert('Tag created successfully');
    } else {
      alert('Tag creation failed');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>
          Tag Name:
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </label>
      </div>
      <button type="submit">Create Tag</button>
    </form>
  );
};

export default CreateTag;

