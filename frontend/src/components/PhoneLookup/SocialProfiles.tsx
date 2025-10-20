// frontend/src/components/PhoneLookup/SocialProfiles.tsx
interface SocialProfilesProps {
  profiles: SocialProfile[];
}

export const SocialProfiles: React.FC<SocialProfilesProps> = ({ profiles }) => {
  const iranPlatforms = {
    'instagram': 'اینستاگرام',
    'telegram': 'تلگرام',
    'whatsapp': 'واتساپ',
    'rubika': 'روبیکا',
    'soroush': 'سروش',
    'eitaa': 'ایتا',
    'gap': 'گپ',
    'twitter': 'توییتر',
    'facebook': 'فیسبوک'
  };

  return (
    <div className="social-profiles">
      {profiles.map((profile, index) => (
        <div key={index} className="profile-card">
          <div className="platform-header">
            <img 
              src={`/platforms/${profile.platform}.png`} 
              alt={iranPlatforms[profile.platform] || profile.platform}
            />
            <span className="platform-name">
              {iranPlatforms[profile.platform] || profile.platform}
            </span>
            <span className={`verification-badge ${profile.verified ? 'verified' : 'unverified'}`}>
              {profile.verified ? 'تایید شده' : 'تایید نشده'}
            </span>
          </div>
          
          <div className="profile-details">
            <div className="username">@{profile.username}</div>
            {profile.displayName && (
              <div className="display-name">{profile.displayName}</div>
            )}
            {profile.followers && (
              <div className="followers">{profile.followers.toLocaleString()} دنبال کننده</div>
            )}
            {profile.lastSeen && (
              <div className="last-seen">
                آخرین فعالیت: {new Date(profile.lastSeen).toLocaleDateString('fa-IR')}
              </div>
            )}
          </div>

          <div className="profile-actions">
            <a href={profile.profileUrl} target="_blank" rel="noopener noreferrer">
              مشاهده پروفایل
            </a>
          </div>
        </div>
      ))}
    </div>
  );
};
