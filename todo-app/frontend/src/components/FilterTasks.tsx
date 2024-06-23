import React, { useState, FormEvent } from 'react';

const FilterTasks: React.FC = () => {
  const [userID, setUserID] = useState<number>(1); // テスト用に固定ユーザーIDを使用
  const [tagID, setTagID] = useState<number | null>(null);
  const [date, setDate] = useState<string>('');

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const params = new URLSearchParams();
    if (userID) params.append('user_id', userID.toString());
    if (tagID) params.append('tag_id', tagID.toString());
    if (date) params.append('date', date);

    const response = await fetch(`/tasks/filter?${params.toString()}`);
    const tasks = await response.json();
    console.log(tasks);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>
          User ID:
          <input
            type="number"
            value={userID}
            onChange={(e) => setUserID(parseInt(e.target.value))}
          />
        </label>
      </div>
      <div>
        <label>
          Tag ID:
          <input
            type="number"
            value={tagID ?? ''}
            onChange={(e) => setTagID(parseInt(e.target.value) || null)}
          />
        </label>
      </div>
      <div>
        <label>
          Date:
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
          />
        </label>
      </div>
      <button type="submit">Filter Tasks</button>
    </form>
  );
};

export default FilterTasks;

