import React, {useEffect} from 'react';

function NginxIgnite() {
  useEffect(
      () => {
        const preloader = document.getElementById('preloader') as HTMLElement
        preloader.remove()
      },
      [],
  )

  return (
      <React.StrictMode>
        <p>Hello there</p>
      </React.StrictMode>
  )
}

export default NginxIgnite
