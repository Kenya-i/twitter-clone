'use client'

import Link from 'next/link'
import { Tweet } from '../types/tweet'

type TweetCardProps = {
  tweet: Tweet
  currentUserId: string | null
  onLike: (tweet: Tweet) => void
  isEditing?: boolean
  editContent?: string
  onEditContentChange?: (value: string) => void
  onEditStart?: (tweet: Tweet) => void
  onEditCancel?: () => void
  onEditSave?: (tweetId: string) => void
  onDelete?: (tweetId: string) => void
}

export default function TweetCard({
  tweet,
  currentUserId,
  onLike,
  isEditing = false,
  editContent = '',
  onEditContentChange,
  onEditStart,
  onEditCancel,
  onEditSave,
  onDelete,
}: TweetCardProps) {
  const canEdit = tweet.user_id === currentUserId && onEditStart && onDelete

  return (
    <div className="border-b border-gray-200 pb-2 flex justify-between items-start">
      <div className="flex-1">
        <Link href={`/users/${tweet.user_id}`} className="text-xs text-blue-500 hover:underline">
          投稿者のプロフィール
        </Link>
        {isEditing ? (
          <div className="mt-1 space-y-2">
            <textarea
              value={editContent}
              onChange={(e) => onEditContentChange?.(e.target.value)}
              className="w-full border border-gray-300 rounded-md px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={2}
            />
            <div className="flex gap-2">
              <button
                onClick={() => onEditSave?.(tweet.id)}
                className="text-xs text-blue-500 hover:text-blue-700"
              >
                保存
              </button>
              <button onClick={onEditCancel} className="text-xs text-gray-500 hover:text-gray-700">
                キャンセル
              </button>
            </div>
          </div>
        ) : (
          <>
            <p className="text-sm mt-1">{tweet.content}</p>
            <p className="text-xs text-gray-400 mt-1">
              {new Date(tweet.created_at).toLocaleString()}
            </p>
            <div className="flex items-center gap-3 mt-1">
              <button
                onClick={() => onLike(tweet)}
                className={`text-xs flex items-center gap-1 ${
                  tweet.liked_by_me ? 'text-pink-500' : 'text-gray-400 hover:text-pink-500'
                }`}
              >
                {tweet.liked_by_me ? '♥' : '♡'} {tweet.like_count}
              </button>
              <Link href={`/tweets/${tweet.id}`} className="text-xs text-blue-500 hover:underline">
                詳細
              </Link>
            </div>
          </>
        )}
      </div>
      {canEdit && !isEditing && (
        <div className="flex gap-2 ml-2">
          <button onClick={() => onEditStart?.(tweet)} className="text-xs text-blue-500 hover:text-blue-700">
            編集
          </button>
          <button onClick={() => onDelete?.(tweet.id)} className="text-xs text-red-500 hover:text-red-700">
            削除
          </button>
        </div>
      )}
    </div>
  )
}
