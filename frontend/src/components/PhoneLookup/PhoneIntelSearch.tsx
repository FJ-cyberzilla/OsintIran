// frontend/src/components/PhoneLookup/PhoneIntelSearch.tsx
import React, { useState } from 'react';
import { useSecureAPI } from '../../hooks/useSecureAPI';
import { encryptPayload } from '../../utils/encryption';

export const PhoneIntelSearch: React.FC = () => {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [searchResults, setSearchResults] = useState<PhoneIntelData | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const { securePost } = useSecureAPI();

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!phoneNumber.startsWith('+98') && !phoneNumber.startsWith('0098')) {
      alert('تنها شماره های ایرانی پشتیبانی می‌شوند'); // Only Iranian numbers supported
      return;
    }

    setIsSearching(true);
    try {
      // Encrypt phone number before sending
      const encryptedData = await encryptPayload({
        phone: phoneNumber,
        timestamp: Date.now(),
        license: await getHardwareLicense()
      });

      const result = await securePost('/api/v1/phone-intel/search', encryptedData);
      setSearchResults(result.data);
    } catch (error) {
      console.error('Search failed:', error);
      alert('خطا در جستجو. لطفا مجددا تلاش کنید.'); // Search error, please try again
    } finally {
      setIsSearching(false);
    }
  };

  return (
    <div className="secure-container">
      <div className="search-header">
        <h1>سامانه هوشمند ردیابی تلفن</h1>
        <p>جستجوی جامع اطلاعات شماره تلفن در فضای مجازی</p>
      </div>

      <form onSubmit={handleSearch} className="search-form">
        <div className="input-group">
          <input
            type="tel"
            value={phoneNumber}
            onChange={(e) => setPhoneNumber(e.target.value)}
            placeholder="+98XXXXXXXXXX یا 0098XXXXXXXXXX"
            pattern="^(\+98|0098)?[0-9]{10}$"
            required
            className="phone-input"
          />
          <button 
            type="submit" 
            disabled={isSearching}
            className="search-btn"
          >
            {isSearching ? 'در حال جستجو...' : 'جستجوی هوشمند'}
          </button>
        </div>
      </form>

      {searchResults && <IntelReport data={searchResults} />}
    </div>
  );
};
